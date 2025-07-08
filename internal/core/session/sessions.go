package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Id        uuid.UUID
	UserId    int
	TokenHash string
	ExpiresAt time.Time
	IssuedAt  time.Time
}

type SessionRepoImpl interface {
	CreateSession(session Session) error
}
