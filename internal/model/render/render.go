package render

import (
	"context"
	"minihost/internal/repository/session"
)

// SessionInfo is basic info about the session (e.g. login state)
// most pages need to know to render templates
type SessionInfo struct {
	LoggedIn bool
	UserUUID string
}

func PopulateSessionInfo(ctx context.Context, s session.Session) (SessionInfo, error) {
	info := SessionInfo{}
	userUUID, err := s.GetAuthenticatedUserUUID(ctx)
	if err != nil {
		return info, err
	}
	if userUUID != nil {
		info.LoggedIn = true
		info.UserUUID = userUUID.String()
	}
	return info, nil
}
