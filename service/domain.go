package service

// User holds information on a user.
type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
