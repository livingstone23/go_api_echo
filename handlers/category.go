package handlers

import (
	"context"
	"encoding/json"
	"go_api_echo/database"
	"go_api_echo/dto"
	"net/http"
	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Category_post save a category in the database
func Category_post(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.CategoryDto

	//Bind the body of the request to the struct
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.Name) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The name is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Save the category in the database
	register := bson.D{{"name", body.Name}, {"slug", slug.Make(body.Name)}}

	database.CategoryCollection.InsertOne(context.TODO(), register)

	answer := map[string]string{
		"State":   "OK",
		"Message": "The category was saved",
	}

	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(answer)
}

// Category_get get all categories
func Category_get(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	findOptions := options.Find()

	cursor, err := database.CategoryCollection.Find(context.TODO(), bson.D{}, findOptions.SetSort(bson.D{{"_id", -1}}))

	if err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the categories",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	var categories []bson.M
	if err = cursor.All(context.TODO(), &categories); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the categories",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	return json.NewEncoder(c.Response()).Encode(categories)
}

// Category_get_by_id get a category by id
func Category_get_by_id(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Get the id from the url
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var category bson.M

	if err := database.CategoryCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&category); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the category",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(category)
}

// Category_put update a category
func Category_put(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.CategoryDto

	//Bind the body of the request to the struct
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Validate the name is not empty
	if len(body.Name) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The name is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Get the id from the url
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	//Update the category in the database
	update := bson.M{"$set": bson.M{"name": body.Name, "slug": slug.Make(body.Name)}}

	if _, err := database.CategoryCollection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, update); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error updating the category",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	answer := map[string]string{
		"State":   "OK",
		"Message": "The category was updated",
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(answer)
}


//Category_delete delete a category
func Category_delete(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Get the id from the url
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	//Validate if the category exists
	if err := database.CategoryCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Err(); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The category does not exist",
		}
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Delete the category in the database
	if _, err := database.CategoryCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID}); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error deleting the category",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	answer := map[string]string{
		"State":   "OK",
		"Message": "The category was deleted",
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(answer)
}