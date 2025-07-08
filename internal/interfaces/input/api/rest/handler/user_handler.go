package handler

import (
	"encoding/json"
	"net/http"

	"taskmgmtsystem/internal/config"
	"taskmgmtsystem/internal/core/users"
	"taskmgmtsystem/pkg/response"
)

type UserHandler struct {
	userService users.UserRepoImpl

	config *config.Config
}

func NewUserHandler(config *config.Config, userService users.UserRepoImpl) UserHandler {
	return UserHandler{
		config:      config,
		userService: userService,
	}
}

func (uh UserHandler) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var user users.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusInternalServerError,
			Error:          err.Error(),
		}
		response.Set()
		return
	}

	insertedUser, err := uh.userService.CreateUser(user)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusInternalServerError,
			Error:          err.Error(),
		}
		response.Set()
		return
	}

	response := response.Response{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Message:        "User created successfully",
		Error:          "none",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Data: insertedUser,
	}
	response.Set()

}
