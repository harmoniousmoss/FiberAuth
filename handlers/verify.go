package handlers

import (
	"context"  // Importing context for managing request context and timeouts
	"net/http" // Importing http for HTTP status codes
	"time"     // Importing time for setting timeouts

	"myfibergotemplate/database" // Importing the custom database package for MongoDB connection
	"myfibergotemplate/models"   // Importing the custom models package for defining the User model

	"github.com/gofiber/fiber/v2"      // Importing Fiber for building the HTTP server
	"go.mongodb.org/mongo-driver/bson" // Importing bson for building MongoDB queries
)

// VerifyEmailHandler handles the email verification process
func VerifyEmailHandler(c *fiber.Ctx) error {
	// Extract the verification token from the query parameters
	token := c.Query("token")
	if token == "" {
		// If the token is missing, return a 400 Bad Request with an error message
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Verification token is missing"})
	}

	// Get a reference to the "users" collection in the MongoDB database
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is cancelled after the operation to prevent resource leaks

	// Define a variable to hold the user information
	var user models.User

	// Search the database for a user with the provided verification token
	err := collection.FindOne(ctx, bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		// If the token is invalid or not found, return a 404 Not Found with an error message
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Invalid verification token"})
	}

	// Create an update document to set the email_status to true and remove the verification_token
	update := bson.M{
		"$set":   bson.M{"email_status": true},     // Set the email status to true (verified)
		"$unset": bson.M{"verification_token": ""}, // Remove the verification token from the database
	}

	// Update the user's document in the database
	_, err = collection.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		// If there's an error during the update, return a 500 Internal Server Error with an error message
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify email"})
	}

	// If everything is successful, return a 200 OK status with a success message
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Email successfully verified"})
}
