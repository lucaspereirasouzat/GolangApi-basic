package functions

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	user "docker.go/src/Models/User"
	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	User user.User
	jwt.StandardClaims
}

// GenerateToken Cria um novo token de usuario
func GenerateToken(user user.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	fmt.Println(user)
	claims := &Token{
		User: user,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//pega a chave do env
	jwtKey := []byte(os.Getenv("API_KEY"))

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ExtractToken faz a extração do Beartoken
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		fmt.Println(strArr[1])
		return strArr[1]
	}
	return ""
}

// VerifyToken faz a verificação do token
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	if len(tokenString) == 0 {
		err := errors.New("an error")
		return nil, err
	}
	// para remover o double coute ""
	newVal := tokenString[1 : len(tokenString)-1]
	// retorna o token autenticado
	token, err := jwt.Parse(newVal, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("API_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
