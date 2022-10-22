package main

import (
	"fmt"
	"lms/database"
	"lms/router"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Start a new fiber app
	nf := fiber.New()

	//connect database
	database.ConnectDB()

	//create router
	router.SetupRoutes(nf)

	nf.Static("/", "./fend/root")
	nf.Static("/api/user/private/", "./fend/private")

	//listen to port 8000
	nf.Listen(":8000")

	t := fmt.Sprint(time.Now().Unix())
	fmt.Println(t)
}
