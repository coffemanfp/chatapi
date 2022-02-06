package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         string
	Nickname   string
	LoggedAt   time.Time
	LastSeenAt time.Time
}

func NewSession(nickname string) (session Session, err error) {
	sessionID, err := HashPassword(nickname + uuid.NewString())
	if err != nil {
		return
	}

	now := time.Now()

	session = Session{
		ID:         sessionID,
		Nickname:   nickname,
		LoggedAt:   now,
		LastSeenAt: now,
	}
	return
}
