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
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteUserHandler deletes a user
func DeleteUserHandler(c *fiber.Ctx) error {
	// Extract the user ID from the URL parameters
	userID := c.Params("id")

	// Convert the userID string to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Extract the authenticated user's ID and role from the context
	authUserID := c.Locals("userID").(string)
	authUserRole := c.Locals("userRole").(string)

	// Ensure that non-administrators can only delete their own data
	if authUserRole != "administrator" && authUserID != userID {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	// Get a reference to the "users" collection in the MongoDB database
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find the user before deleting (this step uses the models.User struct)
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to find user"})
	}

	// Delete the user by ID
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	// Check if a user was deleted
	if result.DeletedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Return a success message
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User deleted successfully"})
}
