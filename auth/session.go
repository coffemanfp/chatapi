package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         string    `json:"id,omitempty"`
	TmpID      string    `json:"tmp_id,omitempty"`
	UserID     int       `json:"nickname,omitempty"`
	LoggedAt   time.Time `json:"logged_at,omitempty"`
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`
	LoggedWith string    `json:"logged_with,omitempty"`
	Actived    bool      `json:"actived,omitempty"`
}

func NewSession(userID int, loggedWith string) (session Session, err error) {
	sessionID, err := HashPassword(fmt.Sprint(userID, uuid.NewString()))
	if err != nil {
		return
	}

	tmpID, err := HashPassword(fmt.Sprint(uuid.NewString(), sessionID))
	if err != nil {
		return
	}

	now := time.Now()

	session = Session{
		ID:         sessionID,
		TmpID:      tmpID,
		UserID:     userID,
		LoggedWith: loggedWith,
		LoggedAt:   now,
		LastSeenAt: now,
		Actived:    true,
	}
	return
}
