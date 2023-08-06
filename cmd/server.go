package cmd

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/michaelcosj/stms/handlers"
	"github.com/michaelcosj/stms/migrations"
	"github.com/michaelcosj/stms/repository"
	"github.com/michaelcosj/stms/router"
)

func Run() error {
	// Setup database
	fmt.Println("INITIALISING DATABASE")

	dbFile := os.Getenv("DB_FILE_PATH")
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("error initialising database: %v", err)
	}
	defer db.Close()

	// Run migrations
	fmt.Println("RUNNING DATABASE MIGRATIONS")

	if err := migrations.RunMigrations(db); err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	// Initialise repository and handlers
	userRepo := repository.InitUserRepo(db)
	handler := handlers.InitHandler(userRepo)

	// Run the router
	port := os.Getenv("PORT")
	if port == "" {
		port = "6969"
	}

	router := router.InitRouter(handler)
	return router.Run(port)
}
