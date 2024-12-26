package handler

import (
	"github.com/keyglee/assessment/lib/todo/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GetTodoRequest struct {
	ID uint `json:"id" param:"id" validate:"required"`
}

func (h *Handler) GetTodo(c echo.Context) error {

	req := GetTodoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, "Invalid request.")
	}

	todo, err := h.TodoRepository.GetByID(req.ID)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(404, "Todo not found.")
	}
	if err != nil {
		h.logger.Errorf("Error getting todo: %v", err)
		return c.JSON(500, "Internal server error.")
	}

	return c.JSON(200, todo)
}

type SearchTodosRequest struct {
	ID          *uint   `param:"id"`
	DisplayName *string `param:"display_name"`
}

func (h *Handler) SearchTodos(c echo.Context) error {
	req := SearchTodosRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, "Invalid request.")
	}

	search := &model.TodoSearchParams{
		ID:          req.ID,
		DisplayName: req.DisplayName,
	}

	todos, err := h.TodoRepository.Search(search)
	h.logger.Infof("Searching todos: %v", todos)
	if err != nil {
		h.logger.Errorf("Error searching todos: %v", err)
		return c.JSON(500, "Internal server error.")
	}

	return c.JSON(200, todos)
}

type CreateTodoRequest struct {
	DisplayName string `json:"display_name"`
}

func (h *Handler) CreateTodo(c echo.Context) error {

	req := CreateTodoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, "Invalid request.")
	}

	toDoItem := &model.Todo{
		DisplayName: req.DisplayName,
	}

	err := h.TodoRepository.Create(toDoItem)
	if err != nil {
		h.logger.Errorf("Error creating todo: %v", err)
		return c.JSON(500, "Internal server error.")
	}

	return c.JSON(201, toDoItem)
}

type UpdateTodoRequest struct {
	ID          uint    `param:"id" validate:"required"`
	DisplayName *string `json:"display_name"`
}

func (h *Handler) UpdateTodo(c echo.Context) error {

	req := UpdateTodoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, "Invalid request.")
	}

	existingTodo, err := h.TodoRepository.GetByID(req.ID)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(404, "Todo not found.")
	}
	if err != nil {
		return c.JSON(500, "Internal server error.")
	}

	update := &model.UpdateTodo{
		DisplayName: req.DisplayName,
	}

	err = h.TodoRepository.Update(existingTodo, update)
	if err != nil {
		h.logger.Errorf("Error updating todo: %v", err)
		return c.JSON(500, "Internal server error.")
	}

	return c.JSON(200, existingTodo)
}

type DeleteTodoRequest struct {
	ID uint `param:"id" validate:"required"`
}

func (h *Handler) DeleteTodo(c echo.Context) error {

	req := DeleteTodoRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, "Invalid request.")
	}

	todo, err := h.TodoRepository.GetByID(req.ID)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(404, "Todo not found.")
	}
	if err != nil {
		h.logger.Errorf("Error getting todo: %v", err)
		return c.JSON(500, "Internal server error.")
	}

	err = h.TodoRepository.Delete(todo)
	if err != nil {
		h.logger.Errorf("Error deleting todo: %v", err)
		return c.JSON(500, "Internal server error.")
	}

	return c.JSON(200, "Todo deleted.")
}
