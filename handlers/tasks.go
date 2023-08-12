package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/models"
)

func (h *handler) AddTask(c echo.Context) error {
	userId := getAuthUserId(c)
	data := make(map[string]interface{})

	t := new(models.Task)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp("error handling request", err))
	}

	// TODO: not sure how i'm to validate these
	tag := strings.ToLower(t.Tag)
	if len(t.Name) < 3 || len(t.Description) < 3 ||
		!(tag == "study" || tag == "work" || tag == "others") {
		data["detail"] = ErrValidatingRequestMsg
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	t.IsCompleted = false
	t.TimeCreated = time.Now()

	// TODO: comeback to this after working on add task on user repo
	taskId, err := h.userRepo.AddTask(userId, *t)
	if err != nil {
		c.Logger().Error(err.Error())
		data["detail"] = "error adding task"
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	t.ID = taskId
	data["task"] = t

	return c.JSON(http.StatusOK, newSuccessResp(data))
}

// TODO: filter and paginate tasks list based on query params
// /users/tasks?is_completed=true&tags=[study,others]&id=10
func (h *handler) GetTasks(c echo.Context) error {
	userId := getAuthUserId(c)
	data := make(map[string]interface{})

	tasks, err := h.userRepo.GetTasks(userId)
	if err != nil {
		c.Logger().Error(err.Error())
		data["detail"] = "error getting task"
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	filter := c.QueryParams()

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
				c.Logger().Error(err.Error())
				data["detail"] = "error parsing query params"
				return c.JSON(http.StatusNotFound, newFailResp(data))
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

	data["tasks"] = filtered_tasks
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) UpdateTask(c echo.Context) error {
	data := make(map[string]interface{})

	newTask := new(models.Task)
	if err := c.Bind(newTask); err != nil {
		return c.JSON(http.StatusInternalServerError, newErrResp("error handling request", err))
	}

	taskIdStr := c.Param("taskId")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		detail := fmt.Sprintf("error parsing taskid %s request: %s", taskIdStr, err.Error())
		data["detail"] = detail
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	if err := h.userRepo.UpdateTask(int64(taskId), *newTask); err != nil {
		data["detail"] = err.Error()
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	data["task"] = newTask
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

func (h *handler) RemoveTask(c echo.Context) error {
	data := make(map[string]interface{})

	taskIdStr := c.Param("taskId")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		data["detail"] = fmt.Sprintf("error parsing taskid %s request: %s", taskIdStr, err.Error())
		return c.JSON(http.StatusBadRequest, newFailResp(data))
	}

	if err := h.userRepo.DeleteTask(int64(taskId)); err != nil {
		data["detail"] = err.Error()
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	data["message"] = "task deleted successfully"
	return c.JSON(http.StatusOK, newSuccessResp(data))
}
