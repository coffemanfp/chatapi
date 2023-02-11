package contact

import (
	"time"

	"github.com/coffemanfp/chat/account"
)

type Contact struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	LastName  string          `json:"last_name"`
	CreatedAt time.Time       `json:"created_at"`
	Account   account.Account `json:"account"`
}
