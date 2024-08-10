package handlers

import (
	"context"
	"net/http"
	"time"

	"myfibergotemplate/database"
	"myfibergotemplate/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUserByIDHandler retrieves a user by their ID
func GetUserByIDHandler(c *fiber.Ctx) error {
	// Extract the user ID from the request URL parameters
	userID := c.Params("id")

	// Convert the userID string to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Extract the user ID from the JWT claims
	jwtUserID := c.Locals("userID").(string)
	jwtUserRole := c.Locals("userRole").(string)

	// Ensure that non-administrators can only access their own data
	if jwtUserRole != "administrator" && jwtUserID != userID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	// Get a reference to the "users" collection in the MongoDB database
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the user by ID
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Return the user data
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User retrieved successfully",
		"user":    user,
	})
}
