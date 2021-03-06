package routes

import (
	"github.com/NazeemNato/tuto/src/handlers"
	"github.com/NazeemNato/tuto/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	// admin
	admin := api.Group("admin")
	admin.Post("register", handlers.Register)
	admin.Post("login", handlers.Login)

	adminAuthenticaed := admin.Use(middlewares.IsAuthenticated)
	adminAuthenticaed.Get("user", handlers.User)
	adminAuthenticaed.Put("user/info", handlers.UpdateInfo)
	adminAuthenticaed.Put("user/password", handlers.UpdatePassword)
	adminAuthenticaed.Post("logout", handlers.Logout)

	adminAuthenticaed.Get("ambassadors", handlers.Ambassador)

	adminAuthenticaed.Get("products", handlers.Products)
	adminAuthenticaed.Post("product", handlers.CreateProduct)
	adminAuthenticaed.Get("product/:id", handlers.GetProduct)
	adminAuthenticaed.Put("product/:id", handlers.UpdateProduct)
	adminAuthenticaed.Delete("product/:id", handlers.DeleteProduct)

	adminAuthenticaed.Get("/user/:id/link", handlers.Link)

	adminAuthenticaed.Get("/orders", handlers.Orders)

	ambassador := api.Group("ambassador")
	ambassador.Post("register", handlers.Register)
	ambassador.Post("login", handlers.Login)
	ambassador.Get("product/frontend", handlers.ProductFrontend)
	ambassador.Get("product/backend", handlers.ProductBackend)

	ambassadorAuthenticaed := ambassador.Use(middlewares.IsAuthenticated)
	ambassadorAuthenticaed.Get("user", handlers.User)
	ambassadorAuthenticaed.Put("user/info", handlers.UpdateInfo)
	ambassadorAuthenticaed.Put("user/password", handlers.UpdatePassword)
	ambassadorAuthenticaed.Post("link", handlers.CreateLink)
	ambassadorAuthenticaed.Get("stats", handlers.Stats)
	ambassadorAuthenticaed.Get("rankings", handlers.Rankings)
	ambassadorAuthenticaed.Post("logout", handlers.Logout)

	checkout := api.Group("checkout")
	checkout.Get("/link/:code", handlers.GetLink)
	checkout.Post("/order", handlers.CreateOrder)
	checkout.Post("/order/confirm", handlers.CompleteOrder)
	
}
