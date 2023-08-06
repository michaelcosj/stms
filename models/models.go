package models

import "time"

type User struct {
	ID         string `json:"id"`
	Tasks      []Task `json:"tasks"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsVerified bool   `json:"is_verified"`
}

type Task struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Priority      uint      `json:"priority"`
	IsCompleted   bool      `json:"is_completed"`
	Tags          []Tag     `json:"tags"`
	Description   string    `json:"description"`
	TimeDue       time.Time `json:"time_due"`
	TimeCreated   time.Time `json:"time_created"`
	TimeCompleted time.Time `json:"time_completed"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
