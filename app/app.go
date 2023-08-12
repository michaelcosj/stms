package app

import (
	"context"

	"github.com/michaelcosj/stms/models"
	"github.com/michaelcosj/stms/repository"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type app struct {
	repo  repository.Repo
	cache *redis.Client
}

type App interface {
	NewUser(username, email, password string) (models.User, error)
	SendVerificationCode(email string) error
	VerifyUser(code string) (models.User, error)
	GetUser(email, password string) (models.User, string, error)

	AddTask(userId int64, t *models.Task) error
	GetTaskByFilters(userId int64, filter map[string][]string) ([]models.Task, error)
	UpdateTask(taskId int64, t *models.Task) error
	DeleteTask(taskId int64) error
}

func InitAppService(repo repository.Repo, cache *redis.Client) App {
	return &app{repo, cache}
}
