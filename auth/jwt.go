package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	Bucket = "AuthBucket"

	OneYear = time.Now().Add(time.Minute * 525600).Unix()
)

type Claims struct {
	Payload []byte `json:"payload"`
	jwt.StandardClaims
}

// CreateToken : new jwt token encrypted with private key
func CreateToken(SigningKey []byte, payload []byte) (string, error) {
	if SigningKey == nil {
		return "", fmt.Errorf("Cannot create token - signing key not provided")
	}

	c := &Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: OneYear,
			Issuer:    "krane-server",
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	// Sign the token with our secret
	tknStr, err := tkn.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return tknStr, nil
}

// ValidateToken : check if token is valid
func ValidateToken(pubKey []byte, tknStr string) bool {
	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return false
	}

	parts := strings.Split(tknStr, ".")
	err = jwt.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return false
	}

	return true
}

// ParseToken : get the contents of a token
func ParseToken(SigningKey []byte, tknStr string) ([]byte, error) {
	token, err := jwt.ParseWithClaims(
		tknStr,    // token
		&Claims{}, // Claims struct
		func(token *jwt.Token) (interface{}, error) {
			return SigningKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims) // Parse into `Claims` struct
	if !ok {
		return nil, errors.New("Could not parse claims")
	}

	// b, _ := json.Marshal(claims.Payload)

	return claims.Payload, nil
}
