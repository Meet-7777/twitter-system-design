package models

type User struct {
	ID           int
	Username     string
	Email        string
	PasswordHash string
	IsVerified   bool
	IsCelebrity  bool
}
