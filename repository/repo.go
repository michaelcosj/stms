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
	users []models.User
}

type UserRepo interface {
	// user management
	NewUser(user models.User) uint
	GetUser(userId uint) (models.User, error)
	GetUserByEmail(userEmail string) (models.User, error)
	UpdateUser(userId uint, user models.User) error
	DeleteUser(userId uint) error

	// task management
	AddTask(userId uint, task models.Task) (uint, error)
	GetTasks(userId uint) ([]models.Task, error)
	UpdateTask(userId uint, taskId uint, task models.Task) error
	DeleteTask(userId uint, taskId uint) error
}

func InitUserRepo() *users {
	return new(users)
}

func (u *users) NewUser(user models.User) uint {
	user.Id = uint(len(u.users)) + 1
	u.users = append(u.users, user)
	return user.Id
}

func (u *users) GetUser(userId uint) (models.User, error) {
	if userId >= uint(len(u.users)) {
		return models.User{}, ErrUserNotFound
	}
	user := u.users[userId-1]

	if user.Id == 0 {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}

func (u *users) GetUserByEmail(userEmail string) (models.User, error) {
	user := models.User{}
	for _, v := range u.users {
		if v.Email == userEmail {
			user = v
		}
	}

	if user.Id == 0 {
		return models.User{}, ErrUserNotFound
	}

	return user, nil
}

func (u *users) UpdateUser(userId uint, user models.User) error {
	if userId >= uint(len(u.users)) || u.users[userId-1].Id == 0 {
		return ErrUserNotFound
	}

	u.users[userId-1] = user

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

	task.Id = uint(len(user.Tasks))
	u.users[userId-1].Tasks = append(user.Tasks, task)

	return task.Id, nil
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

	if taskId >= uint(len(user.Tasks)) || user.Tasks[taskId-1].Id == 0 {
		return ErrTaskNotFound
	}

	user.Tasks[taskId-1] = task
	return u.UpdateUser(user.Id, user)
}

func (u *users) DeleteTask(userId uint, taskId uint) error {
	return u.UpdateTask(userId, taskId, models.Task{})
}
