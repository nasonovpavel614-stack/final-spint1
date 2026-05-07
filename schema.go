package main

import "database/sql"

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS parcel (
	number INTEGER PRIMARY KEY AUTOINCREMENT,
	client INTEGER NOT NULL,
	status TEXT NOT NULL,
	address TEXT NOT NULL,
	created_at TEXT NOT NULL
);`)
	return err
}
