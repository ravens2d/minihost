package model

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
)

// User ...
type User struct {
	UUID         uuid.UUID `db:"uuid"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
}

// NewUser ...
func NewUser(username, email, password string) (*User, error) {
	userUUID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		UUID:         userUUID,
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}, nil
}

// VerifyPassword ...
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))
	return err == nil
}
