package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kshwedha/geodata-go/src/api"
	"github.com/kshwedha/geodata-go/src/common/config"
)

func main() {
	config.Init()

	fmt.Println("Starting server...")

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, userid, token",
		AllowMethods:     "POST, GET, PATCH, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
		ExposeHeaders:    "",
		MaxAge:           3600,
	}))

	app.Post("/login", api.LoginHandler)
	app.Post("/register", api.RegisterHandler)
	app.Post("/upload", api.FileHandler)
	app.Post("/save", api.SaveHandler)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}

	fmt.Println("Running cleanup tasks...")

}
