package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/models"
)

func (h *handler) AddTask(c echo.Context) error {
	userId := getUserIdFromContext(c)

	newTask := new(models.Task)
	if err := c.Bind(newTask); err != nil {
		errMsg := fmt.Sprintf("error handling request: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	if len(newTask.Name) < 3 || len(newTask.Description) < 3 {
		data := map[string]interface{}{"detail": "invalid data", "name": "invalid name", "description": "invalid description"}
		return c.JSON(http.StatusBadRequest, newFailResponse(data))
	}

	newTask.IsCompleted = false
	newTask.TimeCreated = time.Now()

	if _, err := h.userRepo.AddTask(userId, *newTask); err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"task": newTask}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}

func (h *handler) GetTasks(c echo.Context) error {
	userId := getUserIdFromContext(c)

	tasks, err := h.userRepo.GetTasks(userId)
	if err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"tasks": tasks}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}

func (h *handler) UpdateTask(c echo.Context) error {
	userId := getUserIdFromContext(c)

	newTask := new(models.Task)
	if err := c.Bind(newTask); err != nil {
		errMsg := fmt.Sprintf("error handling request: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	taskIdStr := c.Param("taskId")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		detail := fmt.Sprintf("error parsing taskid %s request: %s", taskIdStr, err.Error())
		data := map[string]interface{}{"detail": detail}
		return c.JSON(http.StatusBadRequest, newFailResponse(data))
	}

	if err := h.userRepo.UpdateTask(userId, uint(taskId), *newTask); err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"task": newTask}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}

func (h *handler) RemoveTask(c echo.Context) error {
	userId := getUserIdFromContext(c)

	taskIdStr := c.Param("taskId")
	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		detail := fmt.Sprintf("error parsing taskid %s request: %s", taskIdStr, err.Error())
		data := map[string]interface{}{"detail": detail}
		return c.JSON(http.StatusBadRequest, newFailResponse(data))
	}

	if err := h.userRepo.DeleteTask(userId, uint(taskId)); err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"message": "task " + taskIdStr + " deleted successfully"}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}
