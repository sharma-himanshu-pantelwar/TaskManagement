package users

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserProfile struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

type UserRepoImpl interface {
	CreateUser(user User) (User, error)
}
