package models

import "time"

type User struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Tasks      []Task `json:"tasks"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
}

type Task struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Priority    uint   `json:"priority"`
	IsCompleted bool   `json:"is_completed"`

	Tags        []string `json:"tags"`
	Description string   `json:"description"`

	TimeDue       time.Time `json:"time_due"`
	TimeCreated   time.Time `json:"time_created"`
	TimeCompleted time.Time `json:"time_completed"`
}
