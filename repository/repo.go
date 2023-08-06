package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/michaelcosj/stms/models"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrTaskNotFound = fmt.Errorf("task not found")
)

type users struct {
	users map[string]models.User
}

type UserRepo interface {
	// user management
	NewUser(user models.User) string
	GetUser(userId string) (models.User, error)
	GetUserByEmail(userEmail string) (models.User, error)
	UpdateUser(userId string, user models.User) error
	DeleteUser(userId string) error

	// task management
	AddTask(userId string, task models.Task) (string, error)
	GetTasks(userId string) ([]models.Task, error)
	UpdateTask(userId string, taskId string, task models.Task) error
	DeleteTask(userId string, taskId string) error
}

func InitUserRepo() *users {
	return new(users)
}

func (u *users) NewUser(user models.User) string {
	user.ID = uuid.New().String()
	u.users[user.ID] = user
	return user.ID
}

func (u *users) GetUser(userId string) (models.User, error) {
	for id, user := range u.users {
		if userId == id {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound
}

func (u *users) GetUserByEmail(userEmail string) (models.User, error) {
	for _, user := range u.users {
		if user.Email == userEmail {
			return user, nil
		}
	}

	return models.User{}, ErrUserNotFound

}

func (u *users) UpdateUser(userId string, user models.User) error {
	user, err := u.GetUser(userId)
	if err != nil {
		return err
	}
	u.users[userId] = user

	return nil
}

func (u *users) DeleteUser(userId string) error {
	return u.UpdateUser(userId, models.User{})
}

func (u *users) AddTask(userId string, task models.Task) (string, error) {
	user, err := u.GetUser(userId)
	if err != nil {
		return "", err
	}

	task.ID = uuid.New().String()
	user.Tasks = append(user.Tasks, task)
	u.users[userId] = user

	return task.ID, nil
}

func (u *users) GetTasks(userId string) ([]models.Task, error) {
	user, err := u.GetUser(userId)
	if err != nil {
		return []models.Task{}, err
	}

	return user.Tasks, nil
}

func (u *users) UpdateTask(userId string, taskId string, task models.Task) error {
	user, err := u.GetUser(userId)
	if err != nil {
		return err
	}

	for k, t := range user.Tasks {
		if t.ID == taskId {
			task.ID = taskId
			user.Tasks[k] = task
			u.users[userId] = user
			return nil
		}

	}

	return ErrTaskNotFound
}

func (u *users) DeleteTask(userId string, taskId string) error {
	return u.UpdateTask(userId, taskId, models.Task{})
}
