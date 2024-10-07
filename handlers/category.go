package handlers

import (
	"encoding/json"
	"go_api_echo/database"
	"go_api_echo/dto"
	"net/http"
	"context"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)


//Category_post save a category in the database
func Category_post(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.CategoryDto

	//Bind the body of the request to the struct
	if err:= json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State": "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.Name) == 0 {
		answer := map[string]string{
			"State": "Error",
			"Message": "The name is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Save the category in the database
	register := bson.D{{"name", body.Name}, {"slug", slug.Make(body.Name)} }

	database.CategoryCollection.InsertOne(context.TODO(), register)

	answer := map[string]string{
		"State": "OK",
		"Message": "The category was saved",
	}

	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(answer)
}

//Category_get get all categories
func Category_get(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	findOptions := options.Find()

	cursor, err := database.CategoryCollection.Find(context.TODO(), bson.D{}, findOptions.SetSort(bson.D{{"_id", -1}}))

	if err != nil {
		answer := map[string]string{
			"State": "Error",
			"Message": "Error getting the categories",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	var categories []bson.M
	if err = cursor.All(context.TODO(), &categories); err != nil {
		answer := map[string]string{
			"State": "Error",
			"Message": "Error getting the categories",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	return json.NewEncoder(c.Response()).Encode(categories)
}

