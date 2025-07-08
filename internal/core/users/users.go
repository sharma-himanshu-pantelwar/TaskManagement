package users

import (
	"taskmgmtsystem/internal/config"
	"taskmgmtsystem/internal/core/users/session"
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
}

type UserServiceImpl interface {
	CreateUser(user User) (User, error)
	LoginUser(user User, config *config.Config) (LoginResponse, error)
}
