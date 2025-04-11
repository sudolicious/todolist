package main

import (
	"database/sql"
	"log"
)

func migrate(db *sql.DB) error {
//Check if column created_at exists
	var columnExists bool
	err := db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='tasks' AND column_name='created_at')`).Scan(&columnExists)
	if err != nil {
	return err
}
//add column if not exist, NULL for old tasks
	if !columnExists {
		_, err := db.Exec(`ALTER TABLE tasks ADD COLUMN created_at TIMESTAMP WITH TIME ZONE`)
	if err != nil {
		return err
	}
		log.Println("Migration successful: created_at column added")
	}
	return nil
}
