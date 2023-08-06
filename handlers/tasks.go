package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/michaelcosj/stms/framework"
	"github.com/michaelcosj/stms/models"
)

func (h *handler) AddTask(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	userId := token.Claims.(*framework.CustomClaims).UserID

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
	token := c.Get("user").(*jwt.Token)
	userId := token.Claims.(*framework.CustomClaims).UserID

	tasks, err := h.userRepo.GetTasks(userId)
	if err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"tasks": tasks}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}

func (h *handler) UpdateTask(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	userId := token.Claims.(*framework.CustomClaims).UserID

	newTask := new(models.Task)
	if err := c.Bind(newTask); err != nil {
		errMsg := fmt.Sprintf("error handling request: %s", err.Error())
		return c.JSON(http.StatusInternalServerError, newErrorResponse(errMsg))
	}

	taskId := c.Param("taskId")

	if err := h.userRepo.UpdateTask(userId, taskId, *newTask); err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"task": newTask}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}

func (h *handler) RemoveTask(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	userId := token.Claims.(*framework.CustomClaims).UserID

	taskId := c.Param("taskId")

	if err := h.userRepo.DeleteTask(userId, taskId); err != nil {
		data := map[string]interface{}{"detail": err.Error()}
		return c.JSON(http.StatusNotFound, newFailResponse(data))
	}

	data := map[string]interface{}{"message": "task " + taskId + " deleted successfully"}
	return c.JSON(http.StatusOK, newSuccessResponse(data))
}
