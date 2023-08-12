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
	CheckUserEmailExists(userEmail string) bool
	CheckUserIDExists(userId int64) bool

	// task management
	AddTask(userId int64, task models.Task) (int64, error)
	GetTasks(userId int64) ([]models.Task, error)
	UpdateTask(taskId int64, task models.Task) error
	DeleteTask(taskId int64) error
}

func InitUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db}
}

func (r *userRepo) NewUser(user models.User) (int64, error) {
	res, err := r.db.Exec(insertUserStmt, user.Email, user.Username, user.Password, time.Now().Unix())
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
	var user models.User

	row := r.db.QueryRow(selectUserByIDStmt, userId)
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.IsVerified); err != nil {
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

func (r *userRepo) GetUserByEmail(userEmail string) (models.User, error) {
	var user models.User

	row := r.db.QueryRow(selectUserByEmailStmt, userEmail)
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
	if _, err := r.db.Exec(updateUserStmt, user.Username, user.IsVerified, userId); err != nil {
		return fmt.Errorf("error updating user: %v", err)
	}

	return nil
}

func (r *userRepo) DeleteUser(userId int64) error {
	_, err := r.db.Exec(deleteUserStmt, userId)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

func (r *userRepo) CheckUserEmailExists(userEmail string) bool {
	_, err := r.GetUserByEmail(userEmail)
	return !(err == ErrUserNotFound)
}

func (r *userRepo) CheckUserIDExists(userId int64) bool {
	_, err := r.GetUserByID(userId)
	return !(err == ErrUserNotFound)
}

func (r *userRepo) AddTask(userId int64, t models.Task) (int64, error) {
	res, err := r.db.Exec(insertTaskStmt, t.Name, t.Tag, t.Priority, t.IsCompleted, t.Description, t.TimeDue, t.TimeCreated, userId)
	if err != nil {
		return 0, fmt.Errorf("error inserting task to database: %v", err)
	}

	task_id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return task_id, nil
}

func (r *userRepo) GetTasks(userId int64) ([]models.Task, error) {
	var tasks []models.Task

	row, err := r.db.Query(selectTasksStmt, userId)
	if err != nil {
		return nil, fmt.Errorf("error getting tasks from database: %v", err)
	}

	for row.Next() {
		var t models.Task
		if err := row.Scan(
			&t.ID, &t.Name, &t.Tag, &t.Priority,
			&t.IsCompleted, &t.Description, &t.TimeDue,
			&t.TimeCreated, &t.TimeCompleted,
		); err != nil {
			return nil, fmt.Errorf("error getting task from database: %v", err)
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *userRepo) UpdateTask(id int64, t models.Task) error {
	_, err := r.db.Exec(updateTaskStmt, t.Name, t.Priority, t.IsCompleted, t.Description, t.TimeDue, id)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	return nil
}

// TODO
func (r *userRepo) DeleteTask(taskId int64) error {
	_, err := r.db.Exec(deleteTaskStmt, taskId)
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	return nil
}
