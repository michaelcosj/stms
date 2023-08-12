package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/michaelcosj/stms/models"
)

func (a *app) AddTask(userId int64, t *models.Task) error {
	tag := strings.ToLower(t.Tag)
	if len(t.Name) < 3 || len(t.Description) < 3 ||
		!(tag == "study" || tag == "work" || tag == "others") {
		return fmt.Errorf("invalid task name, description or tag")
	}

	t.IsCompleted = false
	t.TimeCreated = time.Now()

	taskId, err := a.repo.AddTask(userId, *t)
	if err != nil {
		return fmt.Errorf("error adding task to database: %v", err)
	}

	t.ID = taskId
	return err
}

func (a *app) GetTaskByFilters(userId int64, filter map[string][]string) ([]models.Task, error) {
	tasks, err := a.repo.GetTasks(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting tasks from database: %v", err)
	}

	var filtered_tasks []models.Task
	for _, task := range tasks {
		isCompletedFilter := true
		if len(filter["is_completed"]) > 0 {
			isCompletedFilter = strconv.FormatBool(task.IsCompleted) == filter["is_completed"][0]
		}

		priorityFilter := true
		if len(filter["priority"]) > 0 {
			priorityFilter = strconv.FormatBool(task.Priority) == filter["priority"][0]
		}

		idFilter := true
		for _, idStr := range filter["id"] {
			idFilter = false

			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, fmt.Errorf("error parsing id %s", idStr)
			}

			if task.ID == int64(id) {
				idFilter = true
				break
			}
		}

		tagFilter := true
		for _, tag := range filter["tag"] {
			tagFilter = false

			if task.Tag == tag {
				tagFilter = true
				break
			}
		}

		if isCompletedFilter && priorityFilter && idFilter && tagFilter {
			filtered_tasks = append(filtered_tasks, task)
		}
	}

	return filtered_tasks, nil
}

func (a *app) UpdateTask(taskId int64, t *models.Task) error {
	if err := a.repo.UpdateTask(taskId, *t); err != nil {
		return fmt.Errorf("error updating task in database: %v", err)
	}
	return nil
}

func (a *app) DeleteTask(taskId int64) error {
	if err := a.repo.DeleteTask(taskId); err != nil {
		return fmt.Errorf("error removing task from database: %v", err)
	}
	return nil
}
