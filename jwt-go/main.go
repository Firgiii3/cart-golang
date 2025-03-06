package main

import (
	"jwt-go/db"
	"jwt-go/handlers"
	"jwt-go/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	db.InitDB()

	app := fiber.New()

	// Routes
	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)
	app.Get("/user", middleware.AuthMiddleware, handlers.GetUserFromToken)
	app.Delete("/user", middleware.AuthMiddleware, handlers.DeleteUser)
	app.Post("/logout", middleware.AuthMiddleware, handlers.Logout)
	//add cart
	app.Post("/cart", middleware.AuthMiddleware, handlers.AddItemToCart)
	app.Get("/cart/:id", handlers.GetCartItem)
	app.Delete("/cart/:id", middleware.AuthMiddleware, handlers.DeleteCartItem)

	//add product
	app.Post("/product", handlers.AddProduct)
	app.Put("/product/:id", handlers.UpdateProduct)
	app.Delete("/product/:id", handlers.DeleteProduct)

	// Protected Route
	app.Get("/protected", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("You are authorized!")
	})

	log.Fatal(app.Listen(":8080"))
}
