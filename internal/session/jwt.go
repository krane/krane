package session

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/krane/krane/internal/logger"
)

// Claims custom claims for user authentication
type Claims struct {
	Phrase string `json:"phrase"`
	jwt.StandardClaims
}

// CustomClaims custom claims for request
type CustomClaims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

// DecodeJWTToken decodes a jwt token using the signing key
func DecodeJWTToken(signKey string, tknStr string) (jwt.Token, error) {
	tkn, err := jwt.ParseWithClaims(tknStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})

	if err != nil {
		return jwt.Token{}, err
	}

	if !tkn.Valid {
		return jwt.Token{}, errors.New("token is not valid")
	}

	return *tkn, nil
}

// DecodeJWTWithPubKey gets the claims of a jwt auth token
func DecodeJWTWithPubKey(pubKey string, tknStr string) (claims jwt.Claims, err error) {

	// convert ssh format pub key to rsa pub key
	rsaPubKey, err := DecodePublicKey(pubKey)
	if err != nil {
		return
	}

	// validate token signed with private key against rsa public key
	tkn, err := jwt.ParseWithClaims(
		tknStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return rsaPubKey, nil
		},
	)

	if err != nil {
		return
	}

	// verify token is still valid and not expired
	if !tkn.Valid {
		return nil, fmt.Errorf("token not valid %v", tkn.Claims.Valid())
	}

	return tkn.Claims, nil
}

// VerifyAuthTokenWithAuthorizedKeys gets the auth claims from jwt token using an authorized key from server
func VerifyAuthTokenWithAuthorizedKeys(keys []string, tkn string) (claims *Claims) {
	for _, key := range keys {
		c, err := DecodeJWTWithPubKey(key, tkn)
		if err != nil {
			logger.Debugf("unable to decode JWT token with server private key %s", err.Error())
			continue
		}

		// map jwt claims into auth claims
		claims, _ = c.(*Claims)
		break
	}

	return
}

// DecodePublicKey decodes an ssh-rsa string into rsa public key
func DecodePublicKey(str string) (*rsa.PublicKey, error) {
	// comes in as a three part string
	// split into component parts
	tokens := strings.Split(str, " ")

	if len(tokens) < 2 {
		return nil, fmt.Errorf("Invalid key format; must contain at least two fields (keytype data [comment])")
	}

	keyType := tokens[0]
	data, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return nil, err
	}

	format, e, n, err := getRsaValues(data)

	if format != keyType {
		return nil, fmt.Errorf("Key type said %s, but encoded format said %s.  These should match", keyType, format)
	}

	pubKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return pubKey, nil
}

// ParseTokenTypeAndValue returns the type and value of a jwt token
func ParseTokenTypeAndValue(tkn string) (tokenType string, tokenValue string) {
	if tkn == "" {
		logger.Error(errors.New("No token provided"))
		return "", ""
	}
	splitTkn := strings.Split(tkn, " ")
	return splitTkn[0], splitTkn[1]
}

// ParseTokenTypeAndValue returns the session token within the jwt claims
func ParseSessionTokenFromJWTClaims(tkn jwt.Token) (Token, error) {
	claims, ok := tkn.Claims.(*CustomClaims)
	if !ok {
		return Token{}, errors.New("unable to parse the claims for the provided token")
	}

	var sessionTkn Token
	bytes, _ := json.Marshal(claims.Data)
	_ = json.Unmarshal(bytes, &sessionTkn)

	return sessionTkn, nil
}

// IsValidTokenFormat checks if a token is a well formatted Bearer token
func IsValidTokenFormat(tkn string) bool {
	if tkn == "" {
		return false
	}

	// split on the space of the token ex. Bearer XXXXX
	splitTkn := strings.Split(tkn, " ")
	if len(splitTkn) <= 1 {
		return false
	}

	// grab the token type (should be Bearer)
	tknType := splitTkn[0]
	tknValue := splitTkn[1]

	// check if token is a bearer token
	if strings.Compare(tknType, "Bearer") != 0 {
		logger.Infof("Token %s is not a Bearer token", tknType)
		return false
	}

	// if the contents of the token value are empty, then return false
	if tknValue == "" {
		return false
	}

	return true
}

func readLength(data []byte) ([]byte, uint32, error) {
	lBuf := data[0:4]

	buf := bytes.NewBuffer(lBuf)

	var length uint32

	err := binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return nil, 0, err
	}

	return data[4:], length, nil
}

func readBigInt(data []byte, length uint32) ([]byte, *big.Int, error) {
	var bigint = new(big.Int)
	bigint.SetBytes(data[0:length])
	return data[length:], bigint, nil
}

func getRsaValues(data []byte) (format string, e *big.Int, n *big.Int, err error) {
	data, length, err := readLength(data)
	if err != nil {
		return
	}

	format = string(data[0:length])
	data = data[length:]

	data, length, err = readLength(data)
	if err != nil {
		return
	}

	data, e, err = readBigInt(data, length)
	if err != nil {
		return
	}

	data, length, err = readLength(data)
	if err != nil {
		return
	}

	data, n, err = readBigInt(data, length)
	if err != nil {
		return
	}

	return
}
