package database

import (
	"database/sql"
	"minihost/internal/model"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

var (
	db    *sqlx.DB
	cache *redis.Client

	SessionManager *scs.SessionManager
)

const (
	UserUUIDSessionKey = "user_uuid"
)

// TODO: stop abusing globals and init, do some real dependency management lol
func init() {
	var err error

	db, err = sqlx.Connect("sqlite3", "./database.db")
	if err != nil {
		panic(err)
	}

	// TODO: only use one client
	cache = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	SessionManager = scs.New()
	SessionManager.Store = redisstore.New(&redigo.Pool{
		MaxIdle: 10,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", "localhost:6379")
		},
	})
}

// GetUser ...
func GetUser(username string) (*model.User, error) {
	var user model.User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil, nil // no err on not found, just nil model
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser ...
func CreateUser(user *model.User) error {
	_, err := db.Exec("INSERT INTO users (uuid, username, email, password_hash) VALUES (:uuid, :username, :email, :password_hash)", user.UUID, user.Username, user.Email, user.PasswordHash)
	return err
}
