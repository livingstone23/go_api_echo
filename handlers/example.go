package handlers

import (
	"encoding/json"
	"go_api_echo/dto"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func Example_get(c echo.Context) error {
	//return c.String(http.StatusOK, "Hello world from Echo with golang")
	answer := map[string]string{
		"State":   "OK",
		"message": "Hello world from Echo with golang",
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Header().Set("livingstone", "www.livingstone.com")

	//Another way to return a string
	return c.JSON(http.StatusOK, answer)
}

// Receive a header from the request
func Example_get2(c echo.Context) error {
	//return c.String(http.StatusOK, "Hello world from Echo with golang")
	answer := map[string]string{
		"State":   "OK",
		"message": "Hello world from Echo with golang",
		"header":  c.Request().Header.Get("Authorization"),
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Header().Set("livingstone", "www.livingstone.com")

	//Another way to return a string
	return c.JSON(http.StatusOK, answer)
}

func Example_get_with_parameter(c echo.Context) error {

	//The argument must coincide with the parameter in the route
	id := c.Param("id")
	return c.String(http.StatusOK, "Hello world from Echo with golang, id: "+id)

}

func Example_get_with_querystring(c echo.Context) error {

	//The argument must coincide with the parameter in the route
	id := c.QueryParam("id")
	return c.String(http.StatusOK, "Hello world from Echo with golang, id: "+id)

}

func Example_get_json(c echo.Context) error {
	answer := map[string]string{
		"State":   "OK",
		"message": "Hello world from Echo with golang",
	}

	//One way to return json
	//return c.JSON(http.StatusOK, answer)

	//another way return json
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	return json.NewEncoder(c.Response()).Encode(answer)

}

func Example_post(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world from Method post")
}

func Example_post2(c echo.Context) error {

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	var body dto.CategoryDto

	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		awswer := map[string]string{
			"State":   "Error",
			"message": "Error in the body",
		}
		return json.NewEncoder(c.Response()).Encode(awswer)
	}

	awswer := map[string]string{
		"State":   "OK",
		"message": "Hello world from Method post",
		"Name":    body.Name,
	}
	return json.NewEncoder(c.Response()).Encode(awswer)

}

// Function to upload a file
func Upload_file(c echo.Context) error {

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
	var dataFile string = "public/uploads/pictures/" + picture

	dst, err := os.Create(dataFile)
	if err != nil {
		return err
	}

	//Copy the file to the destination
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	anwers := map[string]string{
		"State":   "OK",
		"message": "File " + file.Filename + " uploaded successfully",
		"picture": picture,
	}

	return json.NewEncoder(c.Response()).Encode(anwers)
}

func Example_put(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world from Method put")
}

func Example_delete(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world from Method delete")
}
