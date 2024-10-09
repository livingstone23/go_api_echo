package handlers

import (
	"context"
	"encoding/json"
	"go_api_echo/database"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Fun to post a picture of a product
func ProductPicture_upload(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	file, err := c.FormFile("file")
	if err != nil {
		awswer := map[string]string{
			"State":   "Error",
			"message": "Error in the file",
		}
		return json.NewEncoder(c.Response()).Encode(awswer)
	}

	src, err := file.Open()
	if err != nil {
		awswer := map[string]string{
			"State":   "Error",
			"message": "Error in the file",
		}
		return json.NewEncoder(c.Response()).Encode(awswer)
	}

	//Close the file
	defer src.Close()

	//Create a unique name for the file
	var extension = strings.Split(file.Filename, ".")[1]
	time := strings.Split(time.Now().String(), " ")
	picture := string(time[4][6:14]) + "." + extension
	var dataFile string = "public/uploads/products/" + picture



	dst, err := os.Create(dataFile)
	if err != nil {
		return err
	}

	//Copy the file to the destination
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	//save the register in the database
	objectId, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var result bson.M
	//Validate if exist the product
	if err := database.ProductCollection.FindOne(context.TODO(), bson.M{"_id": objectId}).Decode(&result); err != nil {
		awswer := map[string]string{
			"State":   "Error",
			"message": "The product not exist",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(awswer)
	}

	register := bson.D{{"name",picture}, {"product_id", objectId}}
	database.ProductPictureCollection.InsertOne(context.TODO(), register)


	anwers := map[string]string{
		"State":   "OK",
		"message": "File  register " + file.Filename + " uploaded successfully",
		"picture": picture,
	}

	c.Response().WriteHeader(http.StatusCreated)
	return json.NewEncoder(c.Response()).Encode(anwers)

}

// ProductPicture_get get all pictures of a product
func ProductPicture_get(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Validate if exist the product
	var product bson.M
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := database.ProductCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The product not exist",
		}
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	cursor, err := database.ProductPictureCollection.Find(context.TODO(), bson.M {"product_id": objID })

	if err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the pictures",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	var pictures []bson.M
	if err = cursor.All(context.TODO(), &pictures); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the pictures",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	return json.NewEncoder(c.Response()).Encode(pictures)
}

// ProductPicture_delete delete a picture of a product
func ProductPicture_delete(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//Validate if exist the product
	/*
	var product bson.M
	objID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	if err := database.ProductCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&product); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "The product not exist",
		}
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}
	*/

	//Get the id from the url
	objectID, _ := primitive.ObjectIDFromHex(c.Param("id"))
	var picture bson.M
	if err := database.ProductPictureCollection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&picture); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error getting the picture",
		}

		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	//Delete the picture from the server
	if err := os.Remove("public/uploads/products/" + picture["name"].(string)); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error deleting the picture",
		}
		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}


	//Delete the picture
	if _, err := database.ProductPictureCollection.DeleteOne(context.TODO(), bson.M{"_id": objectID}); err != nil {
		answer := map[string]string{
			"State":   "Error",
			"Message": "Error deleting the picture",
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(c.Response()).Encode(answer)
	}

	answer := map[string]string{
		"State":   "OK",
		"Message": "The picture was deleted",
	}

	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(answer)
}