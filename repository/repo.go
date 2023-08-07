package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/michaelcosj/stms/models"
)

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrTaskNotFound = fmt.Errorf("task not found")
)

type userRepo struct {
	db *sql.DB
}

type UserRepo interface {
	// user management
	NewUser(user models.User) (int64, error)
	GetUserByID(userId int64) (models.User, error)
	GetUserByEmail(userEmail string) (models.User, error)
	UpdateUser(userId int64, user models.User) error
	DeleteUser(userId int64) error

	// task management
	AddTask(userId int64, task models.Task) (int64, error)
	GetTasks(userId int64) ([]models.Task, error)
	UpdateTask(userId int64, taskId int64, task models.Task) error
	DeleteTask(userId int64, taskId int64) error
}

func InitUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db}
}

func (r *userRepo) NewUser(user models.User) (int64, error) {
	res, err := r.db.Exec(insertUserCommand, user.Email, user.Username, user.Password, time.Now().Unix())
	if err != nil {
		return 0, fmt.Errorf("error inserting user to database: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userRepo) GetUserByID(userId int64) (models.User, error) {
	row := r.db.QueryRow(selectUserByIDCommand, userId)

	var user models.User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.IsVerified); err != nil {
		return models.User{}, fmt.Errorf("error getting user from database: %v", err)
	}

	tasks, err := r.GetTasks(user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("error getting user from database: %v", err)
	}

	user.Tasks = tasks
	return user, nil
}

func (r *userRepo) GetUserByEmail(userEmail string) (models.User, error) {
	var user models.User

	row := r.db.QueryRow(selectUserByEmailCommand, userEmail)
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.IsVerified); err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("error getting user from database: %v", err)
	}

	tasks, err := r.GetTasks(user.ID)
	if err != nil {
		return models.User{}, fmt.Errorf("error getting user from database: %v", err)
	}

	if len(tasks) > 0 {
		user.Tasks = tasks
	}

	return user, nil
}

func (r *userRepo) UpdateUser(userId int64, user models.User) error {
	if _, err := r.db.Exec(updateUserCommand, user.Username, user.IsVerified, userId); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func (r *userRepo) DeleteUser(userId int64) error {
	return r.UpdateUser(userId, models.User{})
}

func (u *userRepo) AddTask(userId int64, task models.Task) (int64, error) {
	// user, err := u.GetUser(userId)
	// if err != nil {
	// 	return 0, err
	// }

	// task.ID = int64(len(user.Tasks)) + 1
	// user.Tasks = append(user.Tasks, task)
	// u.users[userId] = user

	// return task.ID, nil
	return 0, nil
}

func (r *userRepo) GetTasks(userId int64) ([]models.Task, error) {
	var tasks []models.Task

	row, err := r.db.Query(selectUserTasksCommand, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting tasks from database: %v", err)
	}

	for row.Next() {
		var t models.Task
		if err := row.Scan(&t.ID, &t.Name, &t.Priority, &t.IsCompleted, &t.Description, &t.TimeDue, &t.TimeCreated, &t.TimeCompleted); err != nil {
			return nil, fmt.Errorf("error getting user from database: %v", err)
		}
		fmt.Println("here 334443343")
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (u *userRepo) UpdateTask(userId int64, taskId int64, task models.Task) error {
	// user, err := u.GetUser(userId)
	// if err != nil {
	// 	return err
	// }

	// for k, t := range user.Tasks {
	// 	if t.ID == taskId {
	// 		task.ID = taskId
	// 		user.Tasks[k] = task
	// 		u.users[userId] = user
	// 		return nil
	// 	}

	// }

	// return ErrTaskNotFound
	return nil
}

func (u *userRepo) DeleteTask(userId int64, taskId int64) error {
	return u.UpdateTask(userId, taskId, models.Task{})
}
