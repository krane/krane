package auth

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// AuthBucket : bucket used for storing auth related key-value data
	AuthBucket = "AuthBucket"

	// SessionsBucket : bucket used for storing session related key-value data
	SessionsBucket = "SessionsBucket"

	// OneYear : unix time for 1 year
	OneYear = time.Now().Add(time.Minute * 525600).Unix()
)

// AuthClaims : custom claims for user authentication
type AuthClaims struct {
	Phrase string `json:"phrase"`
	jwt.StandardClaims
}

// IdentityClaims : custom claims for request identity
type IdentityClaims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

// CreateToken : new jwt token
func CreateToken(SigningKey []byte, data interface{}) (string, error) {
	if SigningKey == nil {
		return "", fmt.Errorf("Cannot create token - signing key not provided")
	}

	customClaims := &IdentityClaims{
		Data: data,
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

// ValidateTokenWithPubKey : check token against public key
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

// ParseAuthToken : get the contents of a jwt token
func ParseAuthToken(signingKey []byte, tknStr string) (phrase string, err error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(signingKey)
	if err != nil {
		return "", err
	}

	tkn, err := jwt.ParseWithClaims(
		tknStr,
		&AuthClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return pubKey, nil
		},
	)
	if err != nil {
		return "", err
	}

	claims, ok := tkn.Claims.(*AuthClaims)
	if ok && tkn.Valid {
		return claims.Phrase, nil
	}

	return "", nil
}
