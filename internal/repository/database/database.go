package database

import (
	"database/sql"
	"errors"
	"minihost/internal/model"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicateUsername = errors.New("username already in use")
	ErrDuplicateEmail    = errors.New("email already in use")
)

type Database interface {
	GetUser(username string) (*model.User, error)
	CreateUser(user *model.User) error
}

type database struct {
	db *sqlx.DB
}

func New() (Database, error) {
	db, err := sqlx.Connect("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	return &database{
		db: db,
	}, nil
}

func (db *database) GetUser(username string) (*model.User, error) {
	var user model.User
	err := db.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil, nil // no err on not found, just nil model
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *database) CreateUser(user *model.User) error {
	_, err := db.db.Exec("INSERT INTO users (uuid, username, email, password_hash) VALUES (:uuid, :username, :email, :password_hash)", user.UUID, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				if sqliteErr.Error() == "UNIQUE constraint failed: users.username" {
					return ErrDuplicateUsername
				}
				if sqliteErr.Error() == "UNIQUE constraint failed: users.email" {
					return ErrDuplicateEmail
				}
			}
		}
	}
	return err
}
