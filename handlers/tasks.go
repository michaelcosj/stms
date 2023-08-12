package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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

	if err := h.app.AddTask(userId, t); err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("error adding task: %v", err))
	}

	data["task"] = t
	return c.JSON(http.StatusOK, newSuccessResp(data))
}

// NOTE: maybe paginate
func (h *handler) GetTasks(c echo.Context) error {
	userId := getAuthUserId(c)
	data := make(map[string]interface{})

	filter := c.QueryParams()

	tasks, err := h.app.GetTaskByFilters(userId, filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("error getting tasks: %v", err))
	}

	data["tasks"] = tasks
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

	if err := h.app.UpdateTask(int64(taskId), newTask); err != nil {
		return c.JSON(http.StatusBadRequest, newErrResp("error updating task", err))
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

	if err := h.app.DeleteTask(int64(taskId)); err != nil {
		data["detail"] = err.Error()
		return c.JSON(http.StatusNotFound, newFailResp(data))
	}

	data["message"] = "task deleted successfully"
	return c.JSON(http.StatusOK, newSuccessResp(data))
}
