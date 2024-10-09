package middleware_custom

import (
	"fmt"
	"context"
	"net/http"
	"os"
	"strings"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go_api_echo/database"
	"github.com/golang-jwt/jwt/v5"


)

//func ValidateJWT
func ValidateJWT(c echo.Context) int {

	errorVar := godotenv.Load()
	if errorVar != nil {
		http.Error(c.Response(), "Error loading .env file", http.StatusUnauthorized)
		return 0
	}

	myKey := []byte(os.Getenv("SECRET_JWT"))
	header := c.Request().Header.Get("Authorization")

	//Check if the header is empty
	if len(header) == 0 {
		http.Error(c.Response(), "The header is required", http.StatusUnauthorized)
		return 0
	}

	//Check if the header has the Bearer token
	splitBearer := strings.Split(header, " ")
	if len(splitBearer) != 2 {
		http.Error(c.Response(), "The header is invalid", http.StatusUnauthorized)
		return 0
	}

	//Check if the token has the 3 parts
	splitToken := strings.Split(splitBearer[1], ".")
	if len(splitToken) != 3 {
		//http.Error(c.Response(), "The token is invalid", http.StatusUnauthorized)
		return 0
	}

	//Check if the token is valid
	tk := strings.TrimSpace(splitBearer[1])
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method:")
		}
		return myKey, nil
	})

	if err != nil {
		http.Error(c.Response(), "The token is invalid", http.StatusUnauthorized)
		return 0
	}

	//Check if the token is valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user bson.M
		//Find the user with the email in the database
		if err := database.UserCollection.FindOne(context.TODO(), bson.D{{"email", claims["email"]}}).Decode(&user); err != nil {
			http.Error(c.Response(), "The user is not in the database", http.StatusUnauthorized)
			return 0
		}
		return 1
	} else {
		http.Error(c.Response(), "The token is invalid", http.StatusUnauthorized)
		return 0
	}
	
}

func handleError(c echo.Context, message string, statusCode int) int {
    http.Error(c.Response(), message, statusCode)
    return 0
}

func ValidateJWT2(c echo.Context) int {
    if err := godotenv.Load(); err != nil {
        return handleError(c, "Error loading .env file", http.StatusUnauthorized)
    }

    myKey := []byte(os.Getenv("SECRET_JWT"))
    header := c.Request().Header.Get("Authorization")

    if len(header) == 0 {
        return handleError(c, "The header is required", http.StatusUnauthorized)
    }

    splitBearer := strings.Split(header, " ")
    if len(splitBearer) != 2 {
        return handleError(c, "The header is invalid", http.StatusUnauthorized)
    }

    splitToken := strings.Split(splitBearer[1], ".")
    if len(splitToken) != 3 {
        return handleError(c, "The token is invalid", http.StatusUnauthorized)
    }

    tk := strings.TrimSpace(splitToken[1])
    token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method")
        }
        return myKey, nil
    })

    if err != nil {
        return handleError(c, "The token is invalid", http.StatusUnauthorized)
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        var user bson.M
        if err := database.UserCollection.FindOne(context.TODO(), bson.D{{"email", claims["email"]}}).Decode(&user); err != nil {
            return handleError(c, "The user is not in the database", http.StatusUnauthorized)
        }
        return 1
    }

    return handleError(c, "The token is invalid", http.StatusUnauthorized)
}