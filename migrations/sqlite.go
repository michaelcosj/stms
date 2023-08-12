package migrations

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	// TODO: priority should be boolean (high and low)
	// TODO: tags should be an enum (study, work, others)
	migrateDbSchema = `
    CREATE TABLE IF NOT EXISTS users (
      user_id         INTEGER   PRIMARY KEY NOT NULL,
      email           TEXT      NOT NULL UNIQUE,
      username        TEXT      NOT NULL,
      password        TEXT      NOT NULL,
      is_verified     BOOLEAN   NOT NULL DEFAULT 0,
      time_created    DATETIME  NOT NULL
    );

    CREATE TABLE IF NOT EXISTS tasks (
      task_id         INTEGER   PRIMARY KEY NOT NULL,
      name            TEXT      NOT NULL,
      TAG             TEXT,
      priority        BOOLEAN   DEFAULT 0, 
      is_completed    BOOLEAN   DEFAULT 0,
      description     TEXT      NOT NULL,
      time_due        DATETIME  NOT NULL,
      time_created    DATETIME  NOT NULL,
      time_completed  DATETIME  NOT NULL DEFAULT 0,
      user_id         INTEGER   NOT NULL REFERENCES users
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
