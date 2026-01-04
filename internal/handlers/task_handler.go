package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/nv-root/task-manager/internal/models"
	"github.com/nv-root/task-manager/internal/services"
	"github.com/nv-root/task-manager/internal/utils"
	"github.com/nv-root/task-manager/internal/validation"
)

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(s *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: s}
}

// disallows fields that are not present in the struct
func DecodeStrict[T any](r io.Reader, v *T) error {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// create task
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	err := DecodeStrict(r.Body, &task)
	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	err = validation.Validate.Struct(task)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		utils.ErrorJSON(w, http.StatusBadRequest, "Validation failed", errs)
		return
	}

	created, err := h.Service.CreateTask(r.Context(), &task)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, "Error creating task", nil)
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, "Task created", created)
}

// get tasks, filter, sort, paginate
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {

	filters := map[string]string{
		"category": "",
		"status":   "",
		"limit":    "",
		"page":     "",
		"sort":     "",
		"order":    "",
		"search":   "",
	}

	for key := range filters {
		if val, ok := r.URL.Query()[key]; ok {
			filters[key] = val[0]
		}
	}

	tasks, err := h.Service.GetTasks(r.Context(), filters)

	if err != nil {
		utils.ErrorJSON(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	utils.ResponseJSON(w, http.StatusOK, "Tasks", struct {
		Count int           `json:"count"`
		Tasks []models.Task `json:"tasks"`
	}{
		Count: len(tasks),
		Tasks: tasks,
	})
}
