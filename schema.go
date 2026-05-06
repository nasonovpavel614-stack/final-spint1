package main

import "database/sql"

// Migrate создаёт таблицу parcel, если её ещё нет.
// Нужна и для main (локальный tracker.db), и для изолированных БД в тестах.
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
