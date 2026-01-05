package routes

import (
	"net/http"

	"github.com/nv-root/task-manager/internal/handlers"
	"github.com/nv-root/task-manager/internal/middleware"
)

func TaskRouter(mux *http.ServeMux, h *handlers.TaskHandler) {
	mux.HandleFunc("POST /tasks", middleware.WithError(h.CreateTask))
	mux.HandleFunc("GET /tasks", middleware.WithError(h.GetTasks))
	mux.HandleFunc("PUT /tasks/{id}", middleware.WithError(h.UpdateTask))
	mux.HandleFunc("DELETE /tasks/{id}", middleware.WithError(h.DeleteTask))
}
