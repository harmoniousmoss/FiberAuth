package handlers

import (
	"context"
	"net/http"
	"time"

	"myfibergotemplate/database"
	"myfibergotemplate/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllUsersHandler retrieves all users from the database
func GetAllUsersHandler(c *fiber.Ctx) error {
	// Get a reference to the "users" collection in the MongoDB database
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Prepare a slice to hold all user records
	var users []models.User

	// Find all user documents in the collection
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		// If there's an error finding the users, return a 500 Internal Server Error with an error message
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode each document into the users slice
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			// If there's an error decoding a document, return a 500 Internal Server Error with an error message
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode user data"})
		}
		users = append(users, user)
	}

	// Check if there were any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cursor iteration error"})
	}

	// Return the list of users in the response
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Users retrieved successfully",
		"users":   users,
	})
}
