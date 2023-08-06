package migrations

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	migrateDbSchema = `
    CREATE TABLE IF NOT EXISTS users (
      id              INTEGER PRIMARY KEY NOT NULL,
      email           TEXT    NOT NULL UNIQUE,
      username        TEXT    NOT NULL,
      password        TEXT    NOT NULL,
      is_verified     BOOLEAN NOT NULL DEFAULT 0,
      time_created    INTEGER NOT NULL
    );

    CREATE TABLE IF NOT EXISTS tasks (
      id              INTEGER NOT NULL,
      name            TEXT    NOT NULL,
      priority        INTEGER DEFAULT 0, 
      is_completed    BOOLEAN DEFAULT 0,
      description     TEXT    NOT NULL,
      time_due        INTEGER NOT NULL,
      time_created    INTEGER NOT NULL,
      time_completed  INTEGER,
      user_id         INTEGER NOT NULL REFERENCES users,
      
      PRIMARY KEY (id, user_id)
    );

    CREATE TABLE IF NOT EXISTS tags (
      id    INTEGER NOT NULL,
      name  TEXT    NOT NULL,
      user_id       INTEGER NOT NULL REFERENCES users,
      
      PRIMARY KEY (id, user_id)
    );

    CREATE TABLE IF NOT EXISTS tag_task_bridge (
      tag_id    INTEGER NOT NULL REFERENCES tags,
      task_id   INTEGER NOT NULL REFERENCES tasks,
      
      PRIMARY KEY (tag_id, task_id)
    );
  `
)

func RunMigrations(db *sql.DB) error {
	if _, err := db.Exec(migrateDbSchema); err != nil {
		return err
	}

	return nil
}

// deletes the database file
func DropDb(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("error dropping database: %v", err)
	}

	return nil
}
