package handlers

import (
	"context"
	"encoding/json"
	"go_api_echo/database"
	"go_api_echo/dto"
	"go_api_echo/jwt"
	"go_api_echo/middleware_custom"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Func Security_register save a user in the database
func Security_register(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.UserDto

	//Bind the body of the request to the struct
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.Email) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The email is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.Password) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The password is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Validate is the email of the user is already in the database
	count, _ := database.UserCollection.CountDocuments(context.TODO(), bson.D{{"email", body.Email}})
	if count > 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The email is already in the database",
		}
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//generate the hash of the password
	cost := 8
	bytes, _ := bcrypt.GenerateFromPassword([]byte(body.Password), cost)

	//Save the user in the database
	register := bson.D{
		{"name", body.Name},
		{"email", body.Email},
		{"telephone", body.Telephone},
		{"password", string(bytes)},
	}

	database.UserCollection.InsertOne(context.TODO(), register)

	answer := map[string]string{
		"State":   "OK",
		"Message": "The user was saved",
	}

	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(answer)
}

// Func Security_login validate the user and return a token
func Security_login(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.LoginDto

	//Bind the body of the request to the struct
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.Email) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The email is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.Password) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The password is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Validate if the user exist
	var user bson.M
	if err := database.UserCollection.FindOne(context.TODO(), bson.M{"email": body.Email}).Decode(&user); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The user not exist",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	passwordbytes := []byte(body.Password)
	passwordBd := []byte(user["password"].(string))

	//Validate the password
	errPassword := bcrypt.CompareHashAndPassword(passwordBd, passwordbytes)

	if errPassword != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The password is incorrect",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	} else {

		//Generate the token
		stringObjectId := user["_id"].(primitive.ObjectID).Hex()
		token, err := jwt.GenerateJWT(user["email"].(string), user["name"].(string), stringObjectId)
		if err != nil {
			answer := map[string]string{
				"State":   "Error",
				"Message": "Error generating the token",
			}

			c.Response().WriteHeader(http.StatusInternalServerError)
			return json.NewEncoder(c.Response()).Encode(answer)
		}

		answer := map[string]string{
			"State":   "OK",
			"user": user["email"].(string),
			"Token":   token,
		}

		c.Response().WriteHeader(http.StatusOK)
		return json.NewEncoder(c.Response()).Encode(answer)

	}

}

// Func to confirm the middleware that confirms the token
func Security_protect(c echo.Context) error {
	
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Validate the token
	if middleware_custom.ValidateJWT(c) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error validating the token",
		}

		c.Response().WriteHeader(http.StatusUnauthorized)
		return json.NewEncoder(c.Response()).Encode(answer)
		
	}

	answer := map[string]string{
		"State":   "OK",
		"Message": "The token is valid",
	}
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(answer)

}