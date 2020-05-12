package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// Bucket : name used for storing auth related key-value data
	Bucket = "AuthBucket"

	// OneYear : unix time for 1 year
	OneYear = time.Now().Add(time.Minute * 525600).Unix()
)

// Claims : custom claims packaged in every token
type Claims struct {
	Payload []byte `json:"payload"`
	jwt.StandardClaims
}

// CreateToken : new jwt token encrypted with a key
func CreateToken(SigningKey []byte, payload []byte) (string, error) {
	if SigningKey == nil {
		return "", fmt.Errorf("Cannot create token - signing key not provided")
	}

	customClaims := &Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: OneYear,
			Issuer:    "krane-server",
		},
	}

	// Declare the unsigned token using RSA HS256 Algorithm for ecryption
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign the token with secret
	signedTkn, err := tkn.SignedString(SigningKey)
	if err != nil {
		return "", err
	}

	return signedTkn, nil
}

// ValidateTokenWithPubKey : check token against publick key
func ValidateTokenWithPubKey(tknStr string) (bool, error) {
	// Read pub key
	pubKey, err := ReadPubKeyFile("")
	if err != nil {
		return false, err
	}

	// Parse public key & verify its PEM encoded
	key, err := jwt.ParseRSAPublicKeyFromPEM(pubKey)
	if err != nil {
		return false, err
	}

	// split on `.` from jwt token
	tknParts := strings.Split(tknStr, ".")

	signingKey := strings.Join(tknParts[0:2], ".")
	signature := tknParts[2]

	err = jwt.SigningMethodRS256.Verify(signingKey, signature, key)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ParseToken : get the contents of a jwt token
func ParseToken(SigningKey []byte, tknStr string) ([]byte, error) {
	token, err := jwt.ParseWithClaims(
		tknStr,
		&Claims{},
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

	return claims.Payload, nil
}
