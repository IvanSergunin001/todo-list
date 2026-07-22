package database

import (
	"database/sql"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id BIGSERIAL PRIMARY KEY,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS dateindex ON scheduler (date);`

func Init(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		db.Close()
		return err
	}

	return nil
}
