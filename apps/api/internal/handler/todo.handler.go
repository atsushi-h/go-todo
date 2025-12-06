package handler

import (
	"net/http"
	"strconv"

	"go-todo/internal/auth"
	"go-todo/internal/model"
	"go-todo/internal/service"

	"github.com/labstack/echo/v4"
)

// TodoのHTTPハンドラー
type TodoHandler struct {
	service *service.TodoService
}

// 新しいTodoHandlerを作成
func NewTodoHandler(service *service.TodoService) *TodoHandler {
	return &TodoHandler{
		service: service,
	}
}

// ListTodos godoc
// @Summary List all todos
// @Description Get all todos from the database
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {array} model.Todo
// @Failure 500 {string} string "Internal server error"
// @Router /todos [get]
func (h *TodoHandler) ListTodos(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	todos, err := h.service.GetAllTodos(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, todos)
}

// GetTodo godoc
// @Summary Get a todo by ID
// @Description Get a single todo by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 200 {object} model.Todo
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Todo not found"
// @Failure 500 {string} string "Internal server error"
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodo(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	id := uint(idInt)
	todo, err := h.service.GetTodoByID(c.Request().Context(), id, userID)
	if err != nil {
		if err == service.ErrTodoNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, todo)
}

// CreateTodo godoc
// @Summary Create a new todo
// @Description Create a new todo with the provided information
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body model.CreateTodoRequest true "Todo to create"
// @Success 201 {object} model.Todo
// @Failure 400 {string} string "Invalid request body or title is required"
// @Failure 500 {string} string "Internal server error"
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var req model.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	if req.Title == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Title is required")
	}

	todo, err := h.service.CreateTodo(c.Request().Context(), req, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusCreated, todo)
}

// UpdateTodo godoc
// @Summary Update a todo
// @Description Update an existing todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Param todo body model.UpdateTodoRequest true "Todo update information"
// @Success 200 {object} model.Todo
// @Failure 400 {string} string "Invalid ID or request body"
// @Failure 404 {string} string "Todo not found"
// @Failure 500 {string} string "Internal server error"
// @Router /todos/{id} [put]
func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	id := uint(idInt)
	var req model.UpdateTodoRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	todo, err := h.service.UpdateTodo(c.Request().Context(), id, userID, req)
	if err != nil {
		if err == service.ErrTodoNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.JSON(http.StatusOK, todo)
}

// DeleteTodo godoc
// @Summary Delete a todo
// @Description Delete a todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Todo not found"
// @Failure 500 {string} string "Internal server error"
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	userID, ok := auth.GetUserIDFromContext(c.Request().Context())
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID")
	}

	id := uint(idInt)

	if err := h.service.DeleteTodo(c.Request().Context(), id, userID); err != nil {
		if err == service.ErrTodoNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return c.NoContent(http.StatusNoContent)
}
