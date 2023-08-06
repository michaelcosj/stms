package models

import "time"

type User struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`

	Tasks []Task `json:"tasks"`
}

type Task struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Priority    int64  `json:"priority"`
	IsCompleted bool   `json:"is_completed"`
	Description string `json:"description"`

	TimeDue       time.Time `json:"time_due"`
	TimeCreated   time.Time `json:"time_created"`
	TimeCompleted time.Time `json:"time_completed"`

	Tags []Tag `json:"tags"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
