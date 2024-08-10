package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/manlikehenryy/blog-backend/controller"
	"github.com/manlikehenryy/blog-backend/middleware"
)

func Setup(app *fiber.App) {
	app.Post("/api/register", controller.Register)
	app.Post("/api/login", controller.Login)

	app.Use(middleware.IsAuthenticated)
	app.Post("/api/post", controller.CreatePost)
	app.Get("/api/post", controller.AllPost)
	app.Get("/api/post/:id", controller.DetailPost)
	app.Put("/api/post/:id", controller.UpdatePost)
	app.Get("/api/user-posts", controller.UsersPost)
	app.Get("/api/logout", controller.Logout)
	app.Delete("/api/post/:id", controller.DeletePost)
	app.Post("/api/upload-image", controller.Upload)
	app.Static("/api/uploads","./uploads")
}