package database

import (
	"context"
	"database/sql"
	"minihost/internal/model"
	"net/http"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/redis/go-redis/v9"
)

const (
	userUUIDSessionKey = "user_uuid"
)

// Repository ...
type Repoistory interface {
	GetUser(username string) (*model.User, error)
	CreateUser(user *model.User) error

	AuthenticateSession(ctx context.Context, user *model.User) error
	GetSessionAuthenticatedUserUUID(ctx context.Context) (*uuid.UUID, error)
	DestroySession(ctx context.Context) error
	SessionLoadAndSave(next http.Handler) http.Handler
}

type repository struct {
	db             *sqlx.DB
	cache          *redis.Client
	sessionManager *scs.SessionManager
}

// NewRepository ...
func NewRepository() (Repoistory, error) {
	db, err := sqlx.Connect("sqlite3", "./database.db")
	if err != nil {
		return nil, err
	}

	cache := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(&redigo.Pool{
		MaxIdle: 10,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", "localhost:6379")
		},
	})

	return &repository{
		db:             db,
		cache:          cache,
		sessionManager: sessionManager,
	}, nil
}

func (r *repository) GetUser(username string) (*model.User, error) {
	var user model.User
	err := r.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err == sql.ErrNoRows {
		return nil, nil // no err on not found, just nil model
	} else if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(user *model.User) error {
	_, err := r.db.Exec("INSERT INTO users (uuid, username, email, password_hash) VALUES (:uuid, :username, :email, :password_hash)", user.UUID, user.Username, user.Email, user.PasswordHash)
	return err
}

func (r *repository) AuthenticateSession(ctx context.Context, user *model.User) error {
	r.sessionManager.Put(ctx, userUUIDSessionKey, user.UUID.String())
	return nil
}

func (r *repository) GetSessionAuthenticatedUserUUID(ctx context.Context) (*uuid.UUID, error) {
	rawUUID := r.sessionManager.GetString(ctx, userUUIDSessionKey)
	if rawUUID == "" {
		return nil, nil // no uuid found if nil
	}
	userUUID, err := uuid.FromString(rawUUID)
	if err != nil {
		return nil, err
	}
	return &userUUID, nil
}

func (r *repository) DestroySession(ctx context.Context) error {
	r.sessionManager.Destroy(ctx)
	return nil
}

func (r *repository) SessionLoadAndSave(next http.Handler) http.Handler {
	return r.sessionManager.LoadAndSave(next)
}
