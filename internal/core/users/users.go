package users

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserProfile struct {
	Username string `json:"username"`
}
