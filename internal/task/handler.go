package task

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sudarshanmg/gotask/internal/auth"
	"github.com/sudarshanmg/gotask/pkg/response"
)

type Handler struct {
	service TaskService
}

func NewHandler(service TaskService) *Handler {
	return &Handler{service: service}
}

func (s *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := auth.GetUserID(r)
	role := auth.GetUserRole(r)

	if userID == 0 {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	task, err := s.service.Create(req, userID, role)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.WriteJSON(w, http.StatusCreated, task)
}

func (s *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	completedStr := r.URL.Query().Get("completed")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	var completed *bool
	if completedStr == "true" {
		b := true
		completed = &b
	} else if completedStr == "false" {
		b := false
		completed = &b
	}

	filter := TaskFilter{
		Completed: completed,
		SortBy:    sort,
		Order:     order,
	}

	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	userID := auth.GetUserID(r)
	role := auth.GetUserRole(r)

	if userID == 0 || role == "" {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tasks, total, totalPages, err := s.service.GetAll(page, limit, filter, userID, role)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}

	w.Header().Set("X-Total-Count", strconv.FormatInt(total, 10))
	w.Header().Set("X-Total-Pages", strconv.Itoa(totalPages))
	w.Header().Set("X-Total-Count", strconv.Itoa(page))
	w.Header().Set("X-Total-Count", strconv.Itoa(limit))

	response.WriteJSON(w, http.StatusOK, tasks)
}

func (s *Handler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid ID format")
		return
	}

	userID := auth.GetUserID(r)
	role := auth.GetUserRole(r)

	task, err := s.service.GetById(id, userID, role)
	if errors.Is(err, ErrInvalidID) {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, ErrNotFound) {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to fetch the task")
		return
	}

	response.WriteJSON(w, http.StatusOK, task)
}

func (s *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid ID format")
		return
	}

	var req UpdateTaskRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userID := auth.GetUserID(r)
	role := auth.GetUserRole(r)

	err = s.service.Update(id, req, userID, role)
	if errors.Is(err, ErrInvalidID) {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, ErrNotFound) {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	log.Printf("Update error: %+v", err)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "task updated successfully"})
}

func (s *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid ID format")
		return
	}

	userID := auth.GetUserID(r)
	role := auth.GetUserRole(r)

	err = s.service.Delete(id, userID, role)

	if errors.Is(err, ErrInvalidID) {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, ErrNotFound) {
		response.WriteError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "task deleted successfully"})
}
