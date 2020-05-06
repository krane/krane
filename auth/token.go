package auth

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userid uint64) (string, error) {
	var err error

	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "biensupernice") //this should be in an env file

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 525600).Unix() // 1 Year
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}
