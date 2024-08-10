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

// ApproveUserHandler approves a user's account by updating their status to "approved"
// but only if their email status is verified
func ApproveUserHandler(c *fiber.Ctx) error {
	// Extract the user ID from the request URL parameters
	userID := c.Params("id")

	// Convert the userID string to a MongoDB ObjectID
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
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

	// Check if the user's email is verified
	if !user.EmailStatus {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Cannot approve user: email not verified"})
	}

	// Update the user's status to "approved"
	update := bson.M{"$set": bson.M{"status": models.Approved, "updated_at": time.Now()}}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to approve user"})
	}

	// Check if the user was found and updated
	if result.MatchedCount == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Return a success message
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User approved successfully"})
}
