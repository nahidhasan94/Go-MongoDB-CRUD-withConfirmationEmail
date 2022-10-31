package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"monGO/config"
	cont "monGO/controller"
)

func main() {
	e := echo.New()

	fmt.Println("Starting Application")
	err := config.InitEnVars() // ?
	if err != nil {
		fmt.Println("[ERROR]: ", err.Error())
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/register", cont.CreateCustomer) // creates user
	e.GET("/getUsers", cont.CustomerList)
	e.DELETE("/delete/:id", cont.Delete)
	e.POST("/update/:id", cont.Update)

	e.Logger.Fatal(e.Start(":3000"))

}
