package cmd

import (
	"fmt"
	"os"

	"github.com/michaelcosj/stms/app"
	"github.com/michaelcosj/stms/framework/cache"
	"github.com/michaelcosj/stms/framework/database"
	"github.com/michaelcosj/stms/handlers"
	"github.com/michaelcosj/stms/migrations"
	"github.com/michaelcosj/stms/repository"
	"github.com/michaelcosj/stms/router"
)

func Run() error {
	// Setup database
	db, err := database.InitDb(os.Getenv("DB_FILE_PATH"))
	if err != nil {
		return fmt.Errorf("error initialising database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := migrations.RunMigrations(db); err != nil {
		return fmt.Errorf("error migrating database: %v", err)
	}

	// setup cache
	cache := cache.InitCache(os.Getenv("REDIS_PORT"))

	// Initialise repository, service and handler
	repo := repository.InitRepo(db)
	service := app.InitAppService(repo, cache)
	handler := handlers.InitHandler(service)

	// Run the router
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "6969"
	}

	router := router.InitRouter(handler)
	return router.Run(port)
}
