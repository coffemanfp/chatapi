package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Session represents the session of a user in a client.
type Session struct {
	ID string `json:"id,omitempty"`

	// TmpID is the just one use ID sent to the client when a new sign in is performed.
	// This ID just must be sent when the client perform the first call to a auth-required route.
	// When the first call is perfomed, the Session.ID must be user for the forward calls.
	TmpID  string `json:"tmp_id,omitempty"`
	UserID int    `json:"user_id,omitempty"`

	// First time that the user has been sign.
	LoggedAt time.Time `json:"logged_at,omitempty"`

	// Last time that the user perform a auth-required route call.
	LastSeenAt time.Time `json:"last_seen_at,omitempty"`

	// Platform which the user has been sign.
	LoggedWith string `json:"logged_with,omitempty"`

	// Session status
	Actived bool `json:"actived,omitempty"`
}

// NewSession initializes a new session instance
//  @param userID int: user id unique identifier.
//  @param loggedWith string: platform which the user has been sign.
//  @return session Session: new Session instance.
//	@return err error: session encryptation error.
func NewSession(userID int, loggedWith string) (session Session, err error) {
	// Generate new encrypted session ID with the userID and a random uuid string.
	sessionID, err := HashPassword(fmt.Sprint(userID, uuid.NewString()))
	if err != nil {
		return
	}

	// Generate a new encrypted temp ID of just one use.
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
