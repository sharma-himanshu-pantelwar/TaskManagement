package handler

import (
	"encoding/json"
	"net/http"

	"taskmgmtsystem/internal/config"
	"taskmgmtsystem/internal/core/users"
	"taskmgmtsystem/pkg/response"
)

type UserHandler struct {
	config      *config.Config
	userService users.UserServiceImpl
}

func NewUserHandler(config *config.Config, userService users.UserServiceImpl) UserHandler {
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
func (uh UserHandler) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req users.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusBadRequest,
			Error:          err.Error(),
		}
		response.Set()
		return
	}

	loginResponse, err := uh.userService.LoginUser(req, uh.config)
	if err != nil {
		response := response.Response{
			ResponseWriter: w,
			StatusCode:     http.StatusInternalServerError,
			Error:          err.Error(),
		}
		response.Set()
		return
	}

	//Set cookie
	accessTokenCookie := http.Cookie{
		Name:     "at",
		Value:    loginResponse.TokenString,
		Expires:  loginResponse.TokenExpire,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Path:     "/",
	}

	sessionCookie := http.Cookie{
		Name:     "sess",
		Value:    loginResponse.Session.Id.String(),
		Expires:  loginResponse.Session.ExpiresAt,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, &accessTokenCookie)
	http.SetCookie(w, &sessionCookie)

	// send response back to client
	response := response.Response{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-user":       loginResponse.FoundUser.Username,
		},
		Message: "Logged in succesfully",
	}

	response.Set()
}
