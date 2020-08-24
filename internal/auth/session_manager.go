package auth

import (
	"encoding/json"
	"errors"

	"github.com/dgrijalva/jwt-go"

	"github.com/biensupernice/krane/internal/storage"
)

var (
	SessionCollection = "Session"
)

// Session : relevant data for authenticated sessions
type Session struct {
	ID        string `json:"id"`
	Principal string `json:"principal"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// Token : for the authenticated session
type SessionToken struct {
	SessionID string `json:"session_id"`
}

// CreateSessionToken : new jwt token
func CreateSessionToken(SigningKey string, sessionTkn SessionToken) (string, error) {
	if SigningKey == "" {
		return "", errors.New("cannot create token - signing key not provided")
	}

	customClaims := &CustomClaims{
		Data: sessionTkn,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: OneYear,
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

func SaveSession(session Session) error {
	if session.ID == "" {
		return errors.New("invalid session")
	}
	bytes, _ := json.Marshal(session)
	return storage.Put(SessionCollection, session.ID, bytes)
}

func GetSessionByID(id string) (Session, error) {
	bytes, err := storage.Get(SessionCollection, id)
	if err != nil {
		return Session{}, err
	}

	var session Session
	_ = json.Unmarshal(bytes, &session)
	return session, nil
}

func GetAllSessions() ([]Session, error) {
	bytes, err := storage.GetAll(SessionCollection)
	if err != nil {
		return make([]Session, 0), err
	}

	sessions := make([]Session, 0)
	for _, session := range bytes {
		var s Session
		_ = json.Unmarshal(session, &s)
		sessions = append(sessions, s)
	}

	return sessions, nil
}
