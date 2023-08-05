package main

import (
	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/handlers"
	"github.com/michaelcosj/stms/repository"
)

func main() {
	e := echo.New()

	userRepo := repository.InitUserRepo()
	handler := handlers.InitHandler(userRepo)

	// User paths
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)

	// Task paths
	e.GET("/users/:userId/tasks", handler.GetTasks)
	e.POST("/users/:userId/tasks", handler.AddTask)
	e.PATCH("/users/:userId/tasks/:taskId", handler.UpdateTask)
	e.DELETE("/users/:userId/tasks/:taskId", handler.RemoveTask)

	e.Logger.Fatal(e.Start(":6969"))
}
