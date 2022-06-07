package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddOrders, downAddOrders)
}

func upAddOrders(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE orders_order (
		    id serial PRIMARY KEY,
			status VARCHAR NOT NULL,
			created_at TIMESTAMP NOT NULL default current_timestamp
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downAddOrders(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE orders_order;")
	if err != nil {
		return err
	}
	return nil
}
