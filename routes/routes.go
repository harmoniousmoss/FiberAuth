package routes

import (
	"myfibergotemplate/handlers"
	"myfibergotemplate/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Signup route
	api.Post("/signup", handlers.SignupHandler)

	// Sign-in route
	api.Post("/signin", handlers.SignInHandler)

	// Email verification route
	api.Get("/verify", handlers.VerifyEmailHandler)

	// Seed admin route
	api.Post("/seed/admin", handlers.SeedAdminHandler)

	// Get all users route - protected by AdminOnlyMiddleware
	api.Get("/users", middleware.AuthMiddleware, middleware.AdminOnlyMiddleware, handlers.GetAllUsersHandler)

	// Approve user route - protected by AdminOnlyMiddleware
	api.Patch("/users/:id/approve", middleware.AuthMiddleware, middleware.AdminOnlyMiddleware, handlers.ApproveUserHandler)

	// Get user by ID route - accessible to the user themselves or administrators
	api.Get("/users/:id", middleware.AuthMiddleware, middleware.OwnDataOrAdminMiddleware, handlers.GetUserByIDHandler)

	// Forgot password route
	api.Post("/forgot-password", handlers.ForgotPasswordHandler)

	// Edit user route - accessible to the user themselves or administrators
	api.Patch("/users/:id", middleware.AuthMiddleware, handlers.EditUserHandler)

	// Delete user route - accessible to the user themselves or administrators
	api.Delete("/users/:id", middleware.AuthMiddleware, handlers.DeleteUserHandler)
}
