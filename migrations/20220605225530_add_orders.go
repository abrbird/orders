package migrations

import (
	"database/sql"
	"fmt"
	"github.com/pressly/goose/v3"
	"gitlab.ozon.dev/zBlur/homework-3/orders/config"
	"log"
	"strings"
)

const (
	minMonthNum      = int64(1)
	maxMonthNum      = int64(12)
	TableName        = "orders_order"
	shardServerNameF = "shard_%d"
	shardTableNameF  = "orders_order_shard_%d"
)

type ShardParameters struct {
	id   int64
	from int64
	to   int64
}

func Split(minV int64, maxV int64, partsNum int64) []ShardParameters {
	shardsParams := make([]ShardParameters, 0)

	rangeV := maxV - minV + 1
	shardsNum := rangeV
	if partsNum < shardsNum {
		shardsNum = partsNum
	}
	step := rangeV / shardsNum

	index := int64(0)
	for int64(len(shardsParams)) < shardsNum {
		from := minV
		if len(shardsParams) > 0 {
			from = shardsParams[len(shardsParams)-1].to + 1
		}
		to := from + step - 1
		if index == (shardsNum - 1) {
			to = maxV
		}
		shardsParams = append(
			shardsParams,
			ShardParameters{
				id:   index,
				from: from,
				to:   to,
			},
		)
		index++
	}

	return shardsParams
}

func init() {
	goose.AddMigration(upAddOrders, downAddOrders)
}

func upAddOrders(tx *sql.Tx) error {
	cfg, err := config.ParseConfig("config/config.yml")
	if err != nil {
		log.Fatal(err)
	}

	mainTableCreation := fmt.Sprintf(`
		CREATE TABLE public.%s (
		    id serial,
			status VARCHAR NOT NULL,
			created_at TIMESTAMP NOT NULL default current_timestamp
		)
		PARTITION BY RANGE (date_part('month', created_at));
	`, TableName)
	shardCreationF := `
		CREATE SERVER IF NOT EXISTS %s FOREIGN DATA WRAPPER postgres_fdw
			OPTIONS (
				dbname '%s',
				host '%s',
				port '%d'
			);
		CREATE USER MAPPING IF NOT EXISTS FOR %s SERVER %s 
			OPTIONS (user '%s', password '%s');
	`
	shardTableCreationF := `
		CREATE FOREIGN TABLE IF NOT EXISTS public.%s
		PARTITION OF public.%s
		FOR VALUES FROM (%d) TO (%d) 
		server %s;
	`
	shardsParams := Split(minMonthNum, maxMonthNum, int64(len(cfg.Database.Shards)))

	queryList := []string{
		mainTableCreation,
		`CREATE EXTENSION IF NOT EXISTS postgres_fdw;`,
		//fmt.Sprintf(`GRANT USAGE ON FOREIGN DATA WRAPPER postgres_fdw to %s;`, cfg.Database.User),
	}

	for i, shardParam := range shardsParams {
		shardConfig := cfg.Database.Shards[i]
		shardServerName := fmt.Sprintf(shardServerNameF, shardParam.id)
		shardTableName := fmt.Sprintf(shardTableNameF, shardParam.id)

		queryList = append(
			queryList,
			fmt.Sprintf(
				shardCreationF,
				shardServerName,
				cfg.Database.DB,
				shardConfig.Host,
				shardConfig.Port,
				cfg.Database.User,
				shardServerName,
				cfg.Database.User,
				cfg.Database.Password,
			),
			fmt.Sprintf(
				shardTableCreationF,
				shardTableName,
				TableName,
				shardParam.from,
				shardParam.to,
				shardServerName,
			),
		)
	}
	queryList = append(
		queryList,
	)

	query := strings.Join(queryList, "")

	_, err = tx.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func downAddOrders(tx *sql.Tx) error {
	_, err := tx.Exec(fmt.Sprintf(`DROP TABLE %s;`, TableName))
	if err != nil {
		return err
	}
	return nil
}
