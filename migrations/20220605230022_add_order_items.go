package migrations

import (
	"database/sql"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAddOrderItems, downAddOrderItems)
}

func upAddOrderItems(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE orders_order_item (
		    id serial PRIMARY KEY,
		    order_id bigint NOT NULL ,
		    volume DECIMAL NOT NULL,
			created_at TIMESTAMP NOT NULL default current_timestamp
			
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downAddOrderItems(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE orders_order_item;")
	if err != nil {
		return err
	}
	return nil
}
