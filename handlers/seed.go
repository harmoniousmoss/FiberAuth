package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"myfibergotemplate/config"
	"myfibergotemplate/database"
	"myfibergotemplate/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdminHandler seeds an admin user into the database
func SeedAdminHandler(c *fiber.Ctx) error {
	// Retrieve the expected seed token from the environment variables
	expectedToken := config.GetEnv("ADMIN_SEED_TOKEN", "")
	if expectedToken == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Admin seed token not configured"})
	}

	// Extract the token from the request header
	providedToken := c.Get("X-Admin-Seed-Token")
	if providedToken != expectedToken {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	// Retrieve admin email and password from the environment variables
	adminEmail := config.GetEnv("ADMIN_EMAIL", "admin@example.com")
	adminPassword := config.GetEnv("ADMIN_PASSWORD", "securepassword")
	adminName := "Admin User"
	adminRole := models.Administrator
	adminStatus := models.Approved

	// Check if the admin user already exists
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": adminEmail}).Decode(&existingUser)
	if err == nil {
		// If the admin user already exists, return a message
		return c.Status(http.StatusConflict).JSON(fiber.Map{"message": "Admin user already exists"})
	}

	// Hash the admin's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Create the admin user object
	adminUser := models.User{
		ID:                 primitive.NewObjectID(),
		MerchantName:       adminName,
		Email:              adminEmail,
		Password:           string(hashedPassword),
		Role:               adminRole,
		Status:             adminStatus,
		EmailStatus:        true, // Email is verified
		PersonInCharge:     adminName,
		PhoneNumber:        "123456789",                 // Example phone number
		Website:            "https://admin.example.com", // Example website
		Address:            "123 Admin Street",          // Example address
		TermsAndConditions: true,                        // Assume admin accepted terms
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Insert the admin user into the database
	_, err = collection.InsertOne(ctx, adminUser)
	if err != nil {
		log.Println("Failed to insert admin user:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to seed admin user"})
	}

	// Return a success message
	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Admin user seeded successfully"})
}
