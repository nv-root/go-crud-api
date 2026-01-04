package routes

import (
	"net/http"

	"github.com/nv-root/task-manager/internal/handlers"
)

func TaskRouter(mux *http.ServeMux, taskHandler *handlers.TaskHandler) {
	mux.HandleFunc("POST /tasks", taskHandler.CreateTask)
	mux.HandleFunc("GET /tasks", taskHandler.GetTasks)
}
