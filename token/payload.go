package token

import (
	"time"

	"github.com/google/uuid"
)

// Payload contains the payload data of the token.
type Payload struct {
	ID        uuid.UUID `json:"id"` // 如果知道 token 洩漏，可以用 id 來註銷它
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}
