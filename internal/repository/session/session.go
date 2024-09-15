package session

import (
	"context"
	"minihost/internal/model"
	"net/http"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gofrs/uuid"
	redigo "github.com/gomodule/redigo/redis"
)

const (
	userUUIDSessionKey = "user_uuid"
)

type Session interface {
	SetAuthenticated(ctx context.Context, user *model.User) error
	GetAuthenticatedUserUUID(ctx context.Context) (*uuid.UUID, error)
	Destroy(ctx context.Context) error
	LoadAndSave(next http.Handler) http.Handler
}

type session struct {
	sessionManager *scs.SessionManager
}

func New() (Session, error) {
	sessionManager := scs.New()
	sessionManager.Store = redisstore.New(&redigo.Pool{
		MaxIdle: 10,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", "localhost:6379")
		},
	})

	return &session{
		sessionManager: sessionManager,
	}, nil
}

func (s *session) SetAuthenticated(ctx context.Context, user *model.User) error {
	s.sessionManager.Put(ctx, userUUIDSessionKey, user.UUID.String())
	return nil
}

func (s *session) GetAuthenticatedUserUUID(ctx context.Context) (*uuid.UUID, error) {
	rawUUID := s.sessionManager.GetString(ctx, userUUIDSessionKey)
	if rawUUID == "" {
		return nil, nil // no uuid found if nil
	}
	userUUID, err := uuid.FromString(rawUUID)
	if err != nil {
		return nil, err
	}
	return &userUUID, nil
}

func (s *session) Destroy(ctx context.Context) error {
	s.sessionManager.Destroy(ctx)
	return nil
}

func (s *session) LoadAndSave(next http.Handler) http.Handler {
	return s.sessionManager.LoadAndSave(next)
}
