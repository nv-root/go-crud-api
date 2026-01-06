package routes

import (
	"net/http"

	"github.com/nv-root/task-manager/internal/handlers"
	"github.com/nv-root/task-manager/internal/middleware"
)

func UserRouter(mux *http.ServeMux, h *handlers.UserHandler) {
	mux.HandleFunc("POST /api/auth/sign-up", middleware.WithError(h.CreateUser))
	mux.HandleFunc("POST /api/auth/login", middleware.WithError(h.LoginUser))

}
