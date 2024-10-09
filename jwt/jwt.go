package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// func to generate the token
func GenerateJWT(email, name, id string) (string, error) {
	errorVA := godotenv.Load()
	if errorVA != nil {
		panic("Error loading .env file")
	}

	myKey := []byte(os.Getenv("SECRET_JWT"))

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":        email,
		"name":         name,
		"created_from": "API-go_api_echo",
		"id":           id,
		"exp":          time.Now().Add(time.Hour * 24).Unix(), //add 24 hours to the token
	})

	tokenString, err := token.SignedString(myKey)
	return tokenString, err
}
