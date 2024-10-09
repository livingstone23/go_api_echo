package main

import (
	"fmt"
	"go_api_echo/handlers"
	"go_api_echo/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	fmt.Println("Starting Program Api Echo...")

	e := echo.New()

	var prefij string = "/api/v1"

	//Static files
	e.Static("/public", "public")

	//Connection to database
	database.ConfirmConnection()

	//Routes
	e.GET(prefij+"/example", handlers.Example_get)
	e.GET(prefij+"/example2", handlers.Example_get2)
	e.GET(prefij+"/example_with_parameters/:id", handlers.Example_get_with_parameter)
	e.GET(prefij+"/example_with_querystring", handlers.Example_get_with_querystring)
	e.GET(prefij+"/example_json", handlers.Example_get_json)
	e.POST(prefij+"/example", handlers.Example_post)
	e.POST(prefij+"/example2", handlers.Example_post2)
	e.PUT(prefij+"/example", handlers.Example_put)
	e.DELETE(prefij+"/example", handlers.Example_delete)
	e.POST(prefij+"/upload_file", handlers.Upload_file)


	//Category routes
	e.POST(prefij+"/category", handlers.Category_post)
	e.GET(prefij+"/category", handlers.Category_get)
	e.GET(prefij+"/category/:id", handlers.Category_get_by_id)
	e.PUT(prefij+"/category/:id", handlers.Category_put)
	e.DELETE(prefij+"/category/:id", handlers.Category_delete)

	//Product routes
	e.POST(prefij+"/product", handlers.Product_post)
	e.GET(prefij+"/product", handlers.Product_get)
	e.GET(prefij+"/products", handlers.Product_get_with_relation)
	e.GET(prefij+"/product/:id", handlers.Product_get_by_id)
	e.PUT(prefij+"/product/:id", handlers.Product_put)
	e.DELETE(prefij+"/product/:id", handlers.Product_delete)

	//Product picture routes
	e.POST(prefij+"/product_picture/:id", handlers.ProductPicture_upload)
	e.GET(prefij+"/product_picture/:id", handlers.ProductPicture_get)
	e.DELETE(prefij+"/product_picture/:id", handlers.ProductPicture_delete)

	//User routes
	e.POST(prefij+"/user", handlers.Security_register)
	e.POST(prefij+"/secure/login", handlers.Security_login)
	e.POST(prefij+"/secure/protect", handlers.Security_protect)


	//Apply cors
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},	
	}))

	errorVars := godotenv.Load()
	if errorVars != nil {
		panic("Error loading .env file")
	}

	e.Logger.Fatal(e.Start(":" + os.Getenv("DB_PORT")))
}
