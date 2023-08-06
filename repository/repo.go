package repository

import (
	"fmt"

	"github.com/michaelcosj/stms/models"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrTaskNotFound = fmt.Errorf("task not found")
)

type users struct {
	users map[uint]models.User
}

type UserRepo interface {
	// user management
	NewUser(user models.User) string
	GetUser(userId uint) (models.User, error)
	GetUserByEmail(userEmail string) (models.User, error)
	UpdateUser(userId uint, user models.User) error
	DeleteUser(userId uint) error

	// task management
	AddTask(userId uint, task models.Task) (string, error)
	GetTasks(userId uint) ([]models.Task, error)
	UpdateTask(userId uint, taskId uint, task models.Task) error
	DeleteTask(userId uint, taskId uint) error
}

func InitUserRepo() *users {
	return new(users)
}

func (u *users) NewUser(user models.User) uint {
	user.ID = uint(len(u.users)) + 1
	u.users[user.ID] = user
	return user.ID
}

func (u *users) GetUser(userId uint) (models.User, error) {
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

func (u *users) UpdateUser(userId uint, user models.User) error {
	user, err := u.GetUser(userId)
	if err != nil {
		return err
	}
	u.users[userId] = user

	return nil
}

func (u *users) DeleteUser(userId uint) error {
	return u.UpdateUser(userId, models.User{})
}

func (u *users) AddTask(userId uint, task models.Task) (uint, error) {
	user, err := u.GetUser(userId)
	if err != nil {
		return 0, err
	}

	task.ID = uint(len(user.Tasks)) + 1
	user.Tasks = append(user.Tasks, task)
	u.users[userId] = user

	return task.ID, nil
}

func (u *users) GetTasks(userId uint) ([]models.Task, error) {
	user, err := u.GetUser(userId)
	if err != nil {
		return []models.Task{}, err
	}

	return user.Tasks, nil
}

func (u *users) UpdateTask(userId uint, taskId uint, task models.Task) error {
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

func (u *users) DeleteTask(userId uint, taskId uint) error {
	return u.UpdateTask(userId, taskId, models.Task{})
}
