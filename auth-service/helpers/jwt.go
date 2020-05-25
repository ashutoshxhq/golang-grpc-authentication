package helpers

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

//CreateToken utility function to create jwt tokens
func CreateToken(u map[string]string) map[string]string {

	//Creating Access Token String
	accessTokenExpirationTime := time.Now().Add(527040 * time.Minute)
	accessClaims := &claims{
		Username: u["username"],
		Role:     u["role"],
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessTokenExpirationTime.Unix(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, _ := accessToken.SignedString(jwtKey)
	token := make(map[string]string)
	token["accessToken"] = accessTokenString
	token["validity"] = time.Now().Add(527040 * time.Minute).String()
	return token
}
