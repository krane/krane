package session

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/krane/krane/internal/auth"
	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/store"
	"github.com/krane/krane/internal/utils"
)

// Session represents an authenticated user session
type Session struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

func (s Session) IsValid() bool {
	if s.ID == "" {
		return false
	}

	if s.User == "" {
		return false
	}

	if s.Token == "" {
		return false
	}

	// TODO: validate expiry date
	return true
}

// CreateSessionToken creates a new jwt token used in a user session instance
func CreateSessionToken(SigningKey string, sessionTkn Token) (string, error) {
	if SigningKey == "" {
		return "", errors.New("cannot create token - signing key not provided")
	}

	customClaims := &auth.CustomClaims{
		Data: sessionTkn,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: utils.OneYear,
			Issuer:    "Krane",
		},
	}

	// Declare the unsigned token using RSA HS256 Algorithm for encryption
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign the token with secret
	signedTkn, err := tkn.SignedString([]byte(SigningKey))
	if err != nil {
		return "", err
	}

	return signedTkn, nil
}

// Save saves a user session into the db
func Save(session Session) error {
	if session.ID == "" {
		return errors.New("invalid session")
	}

	bytes, err := store.Serialize(session)
	if err != nil {
		return err
	}

	return store.Client().Put(constants.SessionsCollectionName, session.ID, bytes)
}

// GetSessionByID returns a user session by id
func GetSessionByID(id string) (Session, error) {
	bytes, err := store.Client().Get(constants.SessionsCollectionName, id)
	if err != nil {
		return Session{}, err
	}

	if bytes == nil {
		return Session{}, fmt.Errorf("session not found")
	}

	var session Session
	err = store.Deserialize(bytes, &session)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

// GetAllSessions returns all user sessions
func GetAllSessions() ([]Session, error) {
	bytes, err := store.Client().GetAll(constants.SessionsCollectionName)
	if err != nil {
		return make([]Session, 0), err
	}

	sessions := make([]Session, 0)
	for _, session := range bytes {
		var s Session
		err := store.Deserialize(session, &s)
		if err != nil {
			return make([]Session, 0), err
		}

		sessions = append(sessions, s)
	}

	return sessions, nil
}
