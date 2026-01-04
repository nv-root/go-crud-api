package handlers

import (
	"encoding/json"
	"fmt"
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
		utils.ErrorJSON(w, 400, "Invalid JSON", err.Error())
		return
	}

	err = validation.Validate.Struct(task)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		utils.ErrorJSON(w, 400, "Validation failed", errs)
		return
	}

	created, err := h.Service.CreateTask(r.Context(), &task)
	if err != nil {
		utils.ErrorJSON(w, 500, "Error creating task", nil)
		return
	}

	utils.ResponseJSON(w, 200, "Task created", created)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Tasks")
}
