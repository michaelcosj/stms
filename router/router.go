package router

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/handlers"
)

type router struct {
	handler handlers.Handler
}

type Router interface {
	Run(port string) error
}

func InitRouter(h handlers.Handler) Router {
	return &router{h}
}

func (r *router) Run(port string) error {
	e := echo.New()

	// Setup logging
	logFile, err := os.Create("log")
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}
	defer logFile.Close()

	logFmt := "time:${time_custom}, method:${method}, uri:${uri}," +
		" status:${status}\nlatency:${latency_human} error:${error}\n\n"

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           logFmt,
		Output:           logFile,
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	// Auth endpoints
	e.POST("/login", r.handler.Login)
	e.POST("/register", r.handler.Register)
	e.POST("/verify", r.handler.VerifyUser)

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

	t.GET("/tasks", r.handler.GetTasks)
	t.POST("/tasks", r.handler.AddTask)
	t.PATCH("/tasks/:taskId", r.handler.UpdateTask)
	t.DELETE("/tasks/:taskId", r.handler.RemoveTask)

	e.Logger.Fatal(e.Start(":" + port))
	return nil
}
