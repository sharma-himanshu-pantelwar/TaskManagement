package usecase

import "taskmgmtsystem/internal/core/users"

type UserService struct {
	userRepo users.UserRepoImpl
}

func NewUserService(userRepo users.UserRepoImpl) users.UserRepoImpl {
	return UserService{userRepo: userRepo}
}

func (us UserService) CreateUser(user users.User) (users.User, error) {
	createdUser, err := us.userRepo.CreateUser(user)
	return createdUser, err
}
