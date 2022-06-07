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
		    id bigint PRIMARY KEY,
		    order_id bigint NOT NULL ,
		    volume DECIMAL NOT NULL,
			created_at TIMESTAMP NOT NULL,
			CONSTRAINT orders_order_item_fk_orders_order 
			    FOREIGN KEY(order_id)
			        REFERENCES orders_order(id) ON DELETE CASCADE
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
