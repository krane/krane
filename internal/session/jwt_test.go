package session

import (
	"testing"

	"github.com/docker/distribution/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSignAndDecodeJWTToken(t *testing.T) {
	signingKey := "something"

	// start by creating a token, then signing it with a key
	tkn := Token{SessionID: uuid.Generate().String()}
	signedTkn, err := CreateSessionToken(signingKey, tkn)
	assert.Nil(t, err)
	assert.NotEqual(t, tkn, signedTkn)

	// decode the token using the same key and parse out the claims
	decodedTkn, err := DecodeJWTToken(signingKey, signedTkn)
	tkn2, err := ParseSessionTokenFromJWTClaims(decodedTkn)

	assert.Nil(t, err)
	assert.True(t, decodedTkn.Valid)
	assert.Equal(t, signedTkn, decodedTkn.Raw)
	assert.Equal(t, "JWT", decodedTkn.Header["typ"])
	assert.Equal(t, "HS256", decodedTkn.Header["alg"])
	assert.Equal(t, "HS256", decodedTkn.Method.Alg())
	assert.Equal(t, tkn.SessionID, tkn2.SessionID)
}

func TestDecodeJWTWithPubKey(t *testing.T) {
	// the signed token was signed using a private key and the matching public key is below
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDbMUO+nJpSXX1mjEl8A5vWlUlHdh7o/Ju8b/9GuER4y/7eGRlk9EuSwCHKaCMVHKEzBSf8XsMJ941ntgrGhuRd62oP+PkgI+gr5SydVoeDfeUTwwmMZsIS689XXE0N81Y1RG+QaaQlcJy88B6utDV0ywD0lZAGsfkhBgqM03f7eXYeFlMLdKUvDDXVfUNjcfeZBwYq9wQREcxPefIWj/Pz8ZFExew/LlUKzhC6NpMTINbnNwBSLC2fn3NJ3nMlVlPEBAcuZZT6ddXYEAEn38Unje6z3EgN1BBxL/ZtWeh2AdbJPLO0tEFUT49lBypY93wmalT95Dop3LQxMbcU1Kms11YZOBDRJP+DYeD3+dx2sEnW6d8jPtuSxfHPgEq4lx7Ueoeq0JxrMcBU1IlelpPerr38tIju6cwMb5rQD5Ds4DKs8XiCfZLcmfjNqsn4QDehUhsqHVXr/RSDmrhliBM6Y/xGt/XwxMPqI8nGzftRWUczNvjxCklNWNDuy4AawXGq4AmzcCL83/N0Pu5DKywA6bc+I1HNb6pZgi5Nj4qOampUDD2fzig9KMadDkaf08rWU10LH9JWMf/6V4gDOzEiDiZ4oARWug6t/I/OMUuVp3V1H5RXcruXiwnK0dj7x6a1kFtsDR/qvK9g9hWMXImheIn7Zb9ah7Pk4PI6g118ew== test@example.com"
	signedTkn := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaHJhc2UiOiJLcmFuZSBhdXRoZW50aWNhdGlvbiByZXF1ZXN0IGlkOiBhOTA1ODAwNC1kNTViLTRlZTUtODMxYy1kY2U4ZjRiYTUwYTkiLCJpYXQiOjE2MTE1MTczNDB9.zlDsZx8hEIOEOsTNn3DsmUqM3JLl2Q8n_0_EPqD1hVDSKfOI6r865dsjpDGd9Q9dp1k3q1n4jHTsbHw56UN6uEgXxvZ47uTY3ArfR6ZMOkj7tnpKUqWLVLDf35136RIm73fwdb8a0rU8Vc_9lXW4kS0J0YX7pGfnpZvc90xx3k4i9DFw4PIzLqx1yH_dS4SEN9RPkshbC5CAi4o0TqANBKiLx-Xf_Oi7oDgEJyhLGXWJvhRfFyzu1ex_WuVpiXo71eVC4wFYgid_0dCk08IYunR-6bPfE2N5SZroEBHNT6nfugdifBQ-uVo83pCrXTWfhUNENILC80fvVkERYoRtEjZ2HGJjuaLzaqMJz-rn5icRO9LKRNWNVWPS8iRckiGQd6fTdmT70LBQpF2YcYA08gIUWuO0Dsn71N1cIk7M_eSutAo6MSRV_uHxSljtYavvWtnE9MnpCIfwe01IXkkBiGb4E3ih-rn2k8Fn09XofpxuaYZjs9dXZbZlh8aDaNwbI2U_S3NCd2gpRSR-_tv-X7y3MoX_HPt9X93eAorK18HyutwvaoCCO7_RS3dxFpEGduySOwZ0vvZv-BW6hAO5ZufyS3eP5dFWppvzVGnEpbVKBXdxlcU4ruPWdDUOjJDoMySLnHupR6ybhwDzAzD9uwUNJrqeYsIs4bibkApfF9A"

	// decode the token using the public key and parse out the claims
	claims, err := DecodeJWTWithPubKey(pubKey, signedTkn)
	assert.Nil(t, err)
	assert.Nil(t, claims.Valid())
}

func TestDecodePublicKey(t *testing.T) {
	pubKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDbMUO+nJpSXX1mjEl8A5vWlUlHdh7o/Ju8b/9GuER4y/7eGRlk9EuSwCHKaCMVHKEzBSf8XsMJ941ntgrGhuRd62oP+PkgI+gr5SydVoeDfeUTwwmMZsIS689XXE0N81Y1RG+QaaQlcJy88B6utDV0ywD0lZAGsfkhBgqM03f7eXYeFlMLdKUvDDXVfUNjcfeZBwYq9wQREcxPefIWj/Pz8ZFExew/LlUKzhC6NpMTINbnNwBSLC2fn3NJ3nMlVlPEBAcuZZT6ddXYEAEn38Unje6z3EgN1BBxL/ZtWeh2AdbJPLO0tEFUT49lBypY93wmalT95Dop3LQxMbcU1Kms11YZOBDRJP+DYeD3+dx2sEnW6d8jPtuSxfHPgEq4lx7Ueoeq0JxrMcBU1IlelpPerr38tIju6cwMb5rQD5Ds4DKs8XiCfZLcmfjNqsn4QDehUhsqHVXr/RSDmrhliBM6Y/xGt/XwxMPqI8nGzftRWUczNvjxCklNWNDuy4AawXGq4AmzcCL83/N0Pu5DKywA6bc+I1HNb6pZgi5Nj4qOampUDD2fzig9KMadDkaf08rWU10LH9JWMf/6V4gDOzEiDiZ4oARWug6t/I/OMUuVp3V1H5RXcruXiwnK0dj7x6a1kFtsDR/qvK9g9hWMXImheIn7Zb9ah7Pk4PI6g118ew== test@example.com"
	rsa, err := DecodePublicKey(pubKey)
	assert.Nil(t, err)
	assert.Equal(t, 512, rsa.Size())
}

func TestParseTokenTypeAndValue(t *testing.T) {
	tkn := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InNlc3Npb25faWQiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIifSwiZXhwIjoxNjQzMDUzMzE4LCJqdGkiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIiLCJpYXQiOjE2MTE1MTczNDAsImlzcyI6IktyYW5lIn0.KeDK7BeZLWfGx2IsdFxCBKj9InRn8X8O8ORL8rc6lKs"
	tknType, tknValue := ParseTokenTypeAndValue(tkn)
	assert.Equal(t, "Bearer", tknType)
	assert.Equal(t, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InNlc3Npb25faWQiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIifSwiZXhwIjoxNjQzMDUzMzE4LCJqdGkiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIiLCJpYXQiOjE2MTE1MTczNDAsImlzcyI6IktyYW5lIn0.KeDK7BeZLWfGx2IsdFxCBKj9InRn8X8O8ORL8rc6lKs", tknValue)
}

func TestIsValidTokenFormat(t *testing.T) {
	assert.False(t, IsValidTokenFormat(""))
	assert.False(t, IsValidTokenFormat("Bearer"))
	assert.False(t, IsValidTokenFormat("Bearer "))
	assert.False(t, IsValidTokenFormat("x"))
	assert.False(t, IsValidTokenFormat("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InNlc3Npb25faWQiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIifSwiZXhwIjoxNjQzMDUzMzE4LCJqdGkiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIiLCJpYXQiOjE2MTE1MTczNDAsImlzcyI6IktyYW5lIn0.KeDK7BeZLWfGx2IsdFxCBKj9InRn8X8O8ORL8rc6lKs"))

	assert.True(t, IsValidTokenFormat("Bearer x"))
	assert.True(t, IsValidTokenFormat("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7InNlc3Npb25faWQiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIifSwiZXhwIjoxNjQzMDUzMzE4LCJqdGkiOiI4ODQ4MTI5MS01YzJjLTQ2YWItODY0ZC04ZjBjYTM2Yjg4ZGIiLCJpYXQiOjE2MTE1MTczNDAsImlzcyI6IktyYW5lIn0.KeDK7BeZLWfGx2IsdFxCBKj9InRn8X8O8ORL8rc6lKs"))
}
