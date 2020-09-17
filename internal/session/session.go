package session

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/biensupernice/krane/internal/auth"
	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

// Session : relevant data for authenticated sessions
type Session struct {
	ID        string `json:"id"`
	Principal string `json:"principal"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// CreateSessionToken : new jwt token
func CreateSessionToken(SigningKey string, sessionTkn Token) (string, error) {
	if SigningKey == "" {
		return "", errors.New("cannot create token - signing key not provided")
	}

	customClaims := &auth.CustomClaims{
		Data: sessionTkn,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: utils.OneYear,
			Issuer:    "krane",
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

func Save(session Session) error {
	if session.ID == "" {
		return errors.New("invalid session")
	}

	bytes, err := store.Serialize(session)
	if err != nil {
		return err
	}

	return store.Instance().Put(collection.Sessions, session.ID, bytes)
}

func GetSessionByID(id string) (Session, error) {
	bytes, err := store.Instance().Get(collection.Sessions, id)
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

func GetAllSessions() ([]Session, error) {
	bytes, err := store.Instance().GetAll(collection.Sessions)
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
