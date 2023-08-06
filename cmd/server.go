package cmd

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/handlers"
	"github.com/michaelcosj/stms/repository"
)

func Run() {
	e := echo.New()

	// Setup logging
	logFile, err := os.Create("log")
	if err != nil {
		log.Fatal("failed to create log file")
	}
	defer logFile.Close()

	logFmt := "time:${time_custom}, method:${method}, uri:${uri}," +
		" status:${status}\nlatency:${latency_human} error:${error}\n\n"

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           logFmt,
		Output:           logFile,
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// Initialise repository and handlers
	userRepo := repository.InitUserRepo()
	handler := handlers.InitHandler(userRepo)

	// Auth endpoints
	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)

	// Task endpoints
	t := e.Group("/users")

	// jwt auth middleware
	jwtMiddlewareCfg := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(framework.CustomClaims)
		},
		SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
	}
	t.Use(echojwt.WithConfig(jwtMiddlewareCfg))

	t.GET("/tasks", handler.GetTasks)
	t.POST("/tasks", handler.AddTask)
	t.PATCH("/tasks/:taskId", handler.UpdateTask)
	t.DELETE("/tasks/:taskId", handler.RemoveTask)

	e.Logger.Fatal(e.Start(":6969"))
}
