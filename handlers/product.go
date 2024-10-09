package handlers

import (
	"context"
	"encoding/json"
	"go_api_echo/database"
	"go_api_echo/dto"
	"net/http"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product_post save a product in the database
func Product_post(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.ProductDto

	//Bind the body of the request to the struct
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Validate the body
	if len(body.Name) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The name of the product required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.CategoryID) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The category is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	CategoryId, _ := primitive.ObjectIDFromHex(body.CategoryID)

	//Save the product in the database
	register := bson.D{
		{"name", body.Name},
		{"price", body.Price},
		{"stock", body.Stock},
		{"description", body.Description},
		{"category_id", CategoryId}}

	database.ProductCollection.InsertOne(context.TODO(), register)

	answer := map[string]string{
		"State":   "OK",
		"Message": "The product was saved",
	}

	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(answer)
}

// Product_get get all products
func Product_get(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	findOptions := options.Find()

	cursor, err := database.ProductCollection.Find(context.TODO(), bson.D{}, findOptions.SetSort(bson.D{{"_id", -1}}))

	if err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the products",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	var results []bson.M

	if err = cursor.All(context.Background(), &results); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the products",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	return c.JSON(http.StatusOK, results)
}

// Product_get_with_relation get all products with relation category
func Product_get_with_relation(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//
	pipeline := []bson.M{
		bson.M{"$match": bson.M{}},
		bson.M{"$lookup": bson.M{"from": "category", "localField": "category_id", "foreignField": "_id", "as": "category"}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}

	cursor, err := database.ProductCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	var results []bson.M

	if err = cursor.All(context.Background(), &results); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the products",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	return c.JSON(http.StatusOK, results)
}

// Product_get_by_id get a product by id
func Product_get_by_id(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Get the id from the url
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	//
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"_id": objectID}},
		bson.M{"$lookup": bson.M{"from": "category", "localField": "category_id", "foreignField": "_id", "as": "category"}},
		bson.M{"$sort": bson.M{"_id": -1}},
	}

	cursor, err := database.ProductCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}

	var results []bson.M

	if err = cursor.All(context.Background(), &results); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the products",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	return c.JSON(http.StatusOK, results[0])
}

// Product_put update a product by id
func Product_put(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.ProductDto

	//Bind the body of the request to the struct
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error in the body",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Validate the body
	if len(body.Name) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The name of the product required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	if len(body.CategoryID) == 0 {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The category is required",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	var result bson.M
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := database.ProductCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&result); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the product",
		}
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Validate the category
	CategoryId, _ := primitive.ObjectIDFromHex(body.CategoryID)

	//Update the product in the database
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{
		"name":        body.Name,
		"price":       body.Price,
		"stock":       body.Stock,
		"description": body.Description,
		"category_id": CategoryId,
	}}

	_, err := database.ProductCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error updating the product",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	answer := map[string]string{
		"State":   "OK",
		"Message": "The product was updated",
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(answer)
}

// Product_delete delete a product
func Product_delete(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Get the id from the url
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))

	//Validate if the product exists
	if err := database.ProductCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Err(); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The product does not exist",
		}
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Delete the product in the database
	_, err := database.ProductCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error deleting the product",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	answer := map[string]string{
		"State":   "OK",
		"Message": "The product was deleted",
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(answer)
}

