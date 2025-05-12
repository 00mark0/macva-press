package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

// Different types of errors returned by the VerifyToken function
var (
	ErrInvalidToken = jwt.ValidationError{Errors: jwt.ValidationErrorSignatureInvalid}
	ErrExpiredToken = jwt.ValidationError{Errors: jwt.ValidationErrorExpired}
)

// Payload contains the payload data of the token, if you want more stuff in the payload add it here
type Payload struct {
	ID            uuid.UUID `json:"id"`
	UserID        string    `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Pfp           string    `json:"pfp"`
	Role          string    `json:"role"`
	EmailVerified bool      `json:"email_verified"`
	Banned        bool      `json:"banned"`
	IsDeleted     bool      `json:"is_deleted"`
	IssuedAt      time.Time `json:"issued_at"`
	ExpiredAt     time.Time `json:"expire_at"`
}

// NewPayload creates a new token payload with a specific username and duration, if you want more stuff in the payload, add it here and make sure it matches your db user table
func NewPayload(userID string, username string, email string, pfp string, role string, emailVerified bool, banned bool, isDeleted bool, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:            tokenID,
		UserID:        userID,
		Username:      username,
		Email:         email,
		Pfp:           pfp,
		Role:          role,
		Banned:        banned,
		IsDeleted:     isDeleted,
		EmailVerified: emailVerified,
		IssuedAt:      time.Now(),
		ExpiredAt:     time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}

	return nil
}

