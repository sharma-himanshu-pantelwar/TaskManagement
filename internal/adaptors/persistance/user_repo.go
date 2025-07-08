package persistance

import (
	"taskmgmtsystem/internal/core/users"
	"taskmgmtsystem/pkg/hashpassword"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) users.UserRepoImpl {
	return UserRepo{db: d}
}

func (u UserRepo) CreateUser(user users.User) (users.User, error) {
	var id int
	query := "INSERT INTO USERS (USERNAME,PASSWORD)VALUES ($1,$2) RETURNING ID"
	hashedPassword, err := hashpassword.HashPassword(user.Password)
	if err != nil {
		return users.User{}, err
	}
	err = u.db.db.QueryRow(query, user.Username, hashedPassword).Scan(&id)
	if err != nil {
		return users.User{}, err
	}
	user.Id = id

	return user, nil
}

func (u UserRepo) FindUserByUsername(username string) (users.User, error) {
	var user users.User
	query := "SELECT ID, USERNAME, PASSWORD FROM USERS WHERE USERNAME=$1;"
	err := u.db.db.QueryRow(query, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return users.User{}, err
	}
	return user, nil
}
