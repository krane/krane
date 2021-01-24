package session

// Token for an authenticated session
type Token struct {
	SessionID string `json:"session_id"` // uuid identifying the session jwt
}

