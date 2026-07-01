package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type Password struct {
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type Claims struct {
	Password string `json:"password"`
	jwt.RegisteredClaims
}

var TokenConst string
var Secret = []byte("my_secret_key")

func signInHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		err := errors.New("method not allowed: must be POST")
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	var password Password

	err := json.NewDecoder(req.Body).Decode(&password)
	if err != nil {
		writeJson(res, http.StatusInternalServerError, err)
		return
	}

	port := os.Getenv("TODO_PASSWORD")
	if port == password.Password {
		jwtToken := jwt.New(jwt.SigningMethodHS256)

		signedToken, err := jwtToken.SignedString(Secret)
		fmt.Println(signedToken)
		if err != nil {
			writeJson(res, http.StatusInternalServerError, err)
			return
		}
		token := Token{
			Token: signedToken,
		}
		TokenConst = signedToken
		jsonData, err := json.Marshal(token)
		if err != nil {
			writeJson(res, http.StatusInternalServerError, err)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(jsonData))
		return
	}
}

func isValidToken(token string, secret []byte) bool {
	claims := Claims{}
	tokenTest, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return false
	}
	if !tokenTest.Valid {
		return false
	}
	if token != TokenConst {
		return false
	}
	return true

}
