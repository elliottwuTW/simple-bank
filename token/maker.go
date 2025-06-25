package token

import "time"

var jwtSigningKey = []byte("simple-bank")

// Maker is an interface for managing tokens.
// If we want to switch different token implementations, it's convenient
// for us to change.
type Maker interface {
	SignToken(username string, duration time.Duration) (string, error)

	VerifyToken(token string) (*Payload, error)
}
