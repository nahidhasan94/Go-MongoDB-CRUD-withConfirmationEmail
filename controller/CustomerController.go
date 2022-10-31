package controller

import (
	"context"
	"fmt"
	"github.com/cakazies/go-mongodb/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"monGO/database"
	"monGO/model"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

//func (u *info) SendEmail(subject, HTMLbody string) error {
//	// sender data
//
//}

func sendMail(customer model.Customer) error {
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PWD")
	to := []string{customer.Email}
	subject := "Account Creation Confirmation\n"
	body := "Hello " + customer.Name + " your Registration has been confirmed.\n Welcome to Klovercloud"
	// smtp - Simple Mail Transfer Protocol
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port
	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, host)

	msg := []byte(subject + body)
	err := smtp.SendMail(address, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func CreateCustomer(c echo.Context) error {
	client, err := database.InitDBConnection()
	if err != nil {
		log.Fatalln("[ERROR] mongo client error: ", err.Error())
	}

	info := model.Customer{}
	resp := model.Resp{}

	err = c.Bind(&info)
	if err != nil {
		log.Printf("Failed processing request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	employeeInfo := client.Database("hr").Collection("employees")

	//checking if user exists

	var checkUser []bson.D

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := employeeInfo.Find(ctx, bson.M{"email": info.Email})
	if err = cursor.All(ctx, &checkUser); err != nil {
		log.Fatal(err)
	}
	if len(checkUser) != 0 {

		resp = model.Resp{Status: false, Message: "This User already exists"}
		return c.JSON(http.StatusBadRequest, resp)
	}

	//

	_, err = employeeInfo.InsertOne(context.Background(), info)

	if err != nil {
		return c.String(http.StatusBadRequest, "BAD REQUEST!")
	}
	if err != nil {
		log.Println("[ERROR] Database insert error: ", err.Error())
	}

	resp = model.Resp{Status: true, Message: "User Created Successfully, Please Check Your Email", Data: info}

	// email

	//go sendMail(info)

	return c.JSON(http.StatusOK, resp)
}

func CustomerList(c echo.Context) error {
	client, err := database.InitDBConnection()
	if err != nil {
		log.Fatalln("[ERROR] mongo client error: ", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //?
	employeeInfo := client.Database("hr").Collection("employees")
	defer cancel()
	cursor, err := employeeInfo.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	var list []bson.M
	if err = cursor.All(ctx, &list); err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, list)
}

func Delete(c echo.Context) error {
	client, err := database.InitDBConnection()
	if err != nil {
		log.Fatalln("[ERROR] mongo client error: ", err.Error())
	}

	params := c.Param("id")
	id, err := primitive.ObjectIDFromHex(params) //?
	utils.LogError(err, "Error Read Params")     //?

	employeeInfo := client.Database("hr").Collection("employees")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := employeeInfo.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Resp(nil, fmt.Sprintf("Error get Per Data : %s", err)))
	}
	fmt.Println(result)
	return c.JSON(http.StatusOK, utils.Resp(nil, fmt.Sprintf("Customer Deleted")))
}

func Update(c echo.Context) error {
	client, err := database.InitDBConnection()
	if err != nil {
		log.Fatalln("[ERROR] mongo client error: ", err.Error())
	}
	info := model.Customer{}

	err = c.Bind(&info)
	if err != nil {
		log.Printf("Failed processing request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	params := c.Param("id")
	id, err := primitive.ObjectIDFromHex(params) //?
	utils.LogError(err, "Error Read Params")

	employeeInfo := client.Database("hr").Collection("employees")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	update := bson.M{"$set": info}

	result, err := employeeInfo.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.Resp(nil, fmt.Sprintf("Error get Per Data : %s", err)))
	}
	fmt.Println(result)
	return c.JSON(http.StatusOK, utils.Resp(nil, fmt.Sprintf("Updated")))
}
