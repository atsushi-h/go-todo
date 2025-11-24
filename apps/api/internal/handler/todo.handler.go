package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-todo/internal/model"
	"go-todo/internal/service"
	"go-todo/internal/util"
)

// パスパラメータを取得（util.GetParamのエイリアス）
func getParam(r *http.Request, name string) string {
	return util.GetParam(r, name)
}

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
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos(r.Context())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
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
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := getParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id := uint(idInt)
	todo, err := h.service.GetTodoByID(r.Context(), id)
	if err != nil {
		if err == service.ErrTodoNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
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
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req model.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo, err := h.service.CreateTodo(r.Context(), req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
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
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := getParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id := uint(idInt)
	var req model.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todo, err := h.service.UpdateTodo(r.Context(), id, req)
	if err != nil {
		if err == service.ErrTodoNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todo)
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
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := getParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id := uint(idInt)

	if err := h.service.DeleteTodo(r.Context(), id); err != nil {
		if err == service.ErrTodoNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
