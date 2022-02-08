package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         string    `json:"id,omitempty"`
	Nickname   string    `json:"nickname,omitempty"`
	LoggedAt   time.Time `json:"logged_at,omitempty"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	LoggedWith string    `json:"logged_with,omitempty"`
}

func NewSession(nickname, loggedWith string) (session Session, err error) {
	sessionID, err := HashPassword(nickname + uuid.NewString())
	if err != nil {
		return
	}

	now := time.Now()

	session = Session{
		ID:         sessionID,
		Nickname:   nickname,
		LoggedWith: loggedWith,
		LoggedAt:   now,
		LastSeenAt: now,
	}
	return
}
