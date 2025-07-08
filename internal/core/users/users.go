package users

import (
	"taskmgmtsystem/internal/config"
	"taskmgmtsystem/internal/core/session"
	"time"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserProfile struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
type LoginResponse struct {
	FoundUser   User
	TokenString string
	TokenExpire time.Time
	Session     session.Session
}
type UserRepoImpl interface {
	CreateUser(user User) (User, error)
	FindUserByUsername(username string) (User, error)
	FindUserById(id int) (UserProfile, error)
}

type UserServiceImpl interface {
	CreateUser(user User) (User, error)
	LoginUser(user User, config *config.Config) (LoginResponse, error)
	GetUserById(userId int) (UserProfile, error)
	LogoutUser(userId int) error
}
