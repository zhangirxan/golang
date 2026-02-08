package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"practice2/internal/models"
	"practice2/internal/storage"
)

const MaxTitleLength = 200


type TaskHandler struct {
	storage *storage.TaskStorage
}

func NewTaskHandler(storage *storage.TaskStorage) *TaskHandler {
	return &TaskHandler{storage: storage}
}


func (h *TaskHandler) HandleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPatch:
		h.handlePatch(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "method not allowed",
		})
	}
}


func (h *TaskHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()


	if idStr := query.Get("id"); idStr != "" {
		h.getTaskByID(w, r, idStr)
		return
	}


	if doneStr := query.Get("done"); doneStr != "" {
		h.getTasksByStatus(w, r, doneStr)
		return
	}


	tasks := h.storage.GetAll()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func (h *TaskHandler) getTaskByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid id",
		})
		return
	}

	task, exists := h.storage.GetByID(id)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "task not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}


func (h *TaskHandler) getTasksByStatus(w http.ResponseWriter, r *http.Request, doneStr string) {
	done, err := strconv.ParseBool(doneStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid done parameter",
		})
		return
	}

	tasks := h.storage.GetByStatus(done)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}


func (h *TaskHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTaskRequest
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}


	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid title",
		})
		return
	}

	if len(req.Title) > MaxTitleLength {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "title too long (max 200 characters)",
		})
		return
	}


	task := h.storage.Create(req.Title)
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}


func (h *TaskHandler) handlePatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "id parameter is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid id",
		})
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid request body",
		})
		return
	}


	updated := h.storage.Update(id, req.Done)
	if !updated {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "task not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.UpdateResponse{
		Updated: true,
	})
}

func (h *TaskHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "id parameter is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "invalid id",
		})
		return
	}

	deleted := h.storage.Delete(id)
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error: "task not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{
		"deleted": true,
	})
}
