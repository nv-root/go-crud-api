package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nv-root/task-manager/internal/models"
	"github.com/nv-root/task-manager/internal/services"
	"github.com/nv-root/task-manager/internal/utils"
	"github.com/nv-root/task-manager/internal/validation"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		Service: service,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var user models.CreateUserRequest

	err := DecodeStrict(r.Body, &user)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(user)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	created, err := h.Service.CreateUser(r.Context(), &user)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusCreated, "User created", created)
	return nil
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) error {
	start := time.Now()
	defer func() {
		fmt.Println("Handler took:", time.Since(start))
	}()

	var creds models.Credentials

	err := DecodeStrict(r.Body, &creds)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(creds)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	data, err := h.Service.LoginUser(r.Context(), &creds)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Loggedin", data)
	return nil
}

func (h *UserHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) error {
	var user models.CreateUserRequest

	err := DecodeStrict(r.Body, &user)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = h.Service.ForgotPassword(r.Context(), user.Email)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Reset password email has been sent to your email", nil)
	return nil
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) error {
	var req models.UpdatePasswordRequest

	token := r.URL.Query().Get("token")
	if token == "" {
		return utils.BadRequest("Missing reset password token", nil)
	}

	err := DecodeStrict(r.Body, &req)
	if err != nil {
		return utils.BadRequest("Invalid JSON", nil)
	}

	err = validation.Validate.Struct(req)
	if err != nil {
		errs := utils.FormatValidationErrors(err)
		return utils.BadRequest("Validation Failed", errs)
	}

	err = h.Service.ResetPassword(r.Context(), token, &req)
	if err != nil {
		return err
	}

	utils.ResponseJSON(w, http.StatusOK, "Password updated", nil)
	return nil
}
