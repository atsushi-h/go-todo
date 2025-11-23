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

// 全てのTodoを取得
func (h *TodoHandler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(todos)
}

// 指定されたIDのTodoを取得
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	idStr := getParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id := uint(idInt)
	todo, err := h.service.GetTodoByID(id)
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

// 新しいTodoを作成
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

	todo, err := h.service.CreateTodo(req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// 既存のTodoを更新
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

	todo, err := h.service.UpdateTodo(id, req)
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

// 指定されたIDのTodoを削除
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	idStr := getParam(r, "id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil || idInt < 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	id := uint(idInt)

	if err := h.service.DeleteTodo(id); err != nil {
		if err == service.ErrTodoNotFound {
			http.Error(w, "Todo not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
