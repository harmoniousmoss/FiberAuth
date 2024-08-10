package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"myfibergotemplate/config"
	"myfibergotemplate/database"
	"myfibergotemplate/libs"
	"myfibergotemplate/models"
	"myfibergotemplate/utils" // Import the utils package

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// SignupHandler handles the user signup process
func SignupHandler(c *fiber.Ctx) error {
	user := new(models.User) // Create a new instance of the User model

	// Parse the incoming JSON request body into the user struct
	if err := c.BodyParser(user); err != nil {
		// If there's an error in parsing, return a 400 Bad Request with an error message
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Set the default values for the user object
	user.ID = primitive.NewObjectID() // Generate a new MongoDB ObjectID
	user.Status = models.Pending      // Set the default status to "pending"
	user.Role = models.Merchant       // Set the default role to "merchant"
	user.EmailStatus = false          // Set the default email verification status to "false"
	user.CreatedAt = time.Now()       // Set the current time as the creation time
	user.UpdatedAt = time.Now()       // Set the current time as the last updated time

	// Hash the user's password using bcrypt for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		// If there's an error hashing the password, return a 500 Internal Server Error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(hashedPassword) // Store the hashed password in the user object

	// Generate a random token for email verification
	verificationToken := utils.GenerateVerificationToken() // Use the utility function
	user.VerificationToken = verificationToken             // Store the verification token in the user object

	// Get a reference to the "users" collection in the MongoDB database (using the correct database name)
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	// Insert the new user into the "users" collection
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		// If there's an error inserting the user, return a 500 Internal Server Error
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	// Prepare the email content for verification
	backendURL := config.GetEnv("BACKEND_URL", "http://localhost:8080")       // Get the backend URL from the environment variables
	verificationLink := backendURL + "/api/verify?token=" + verificationToken // Create the verification link with the token

	emailSubject := "Email registration verification - TalentDev ID"
	emailBody := buildVerificationEmail(user.MerchantName, verificationLink)

	// Send email to the user for email verification
	err = libs.SendEmail([]string{user.Email}, emailSubject, emailBody)
	if err != nil {
		// If there's an error sending the email, log the error (but don't return it to the user)
		log.Println("Failed to send verification email to user:", err)
	}

	// Return a 201 Created status with a success message and the user's data
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "New user signup successfully created",
		"user": fiber.Map{
			"id":                   user.ID.Hex(),           // Convert the ObjectID to a string
			"merchant_name":        user.MerchantName,       // Return the merchant name
			"status":               user.Status,             // Return the user's status
			"email":                user.Email,              // Return the user's email
			"email_status":         user.EmailStatus,        // Return the email verification status
			"role":                 user.Role,               // Return the user's role
			"person_in_charge":     user.PersonInCharge,     // Return the person in charge
			"phone_number":         user.PhoneNumber,        // Return the phone number
			"website":              user.Website,            // Return the website
			"address":              user.Address,            // Return the address
			"terms_and_conditions": user.TermsAndConditions, // Return whether the user accepted the terms and conditions
			"created_at":           user.CreatedAt,          // Return the creation time
			"updated_at":           user.UpdatedAt,          // Return the last updated time
		},
	})
}

// buildVerificationEmail constructs the email body using the provided merchant name and verification link
func buildVerificationEmail(merchantName, verificationLink string) string {
	emailTemplate := `
		<p>Dear ${merchantName},</p>
		<p>Welcome to TalentDev!</p>
		<p>Thank you for signing up on our platform. To complete your registration, please verify your email address by clicking the link below:</p>
		<p><a href="${verificationLink}">Click here to verify your email address</a></p>
		<p>If you did not sign up for a TalentDev account, please ignore this email.</p>
		<p>Kind regards,<br>The TalentDev Team</p>
	`

	// Replace placeholders with actual values
	emailTemplate = strings.ReplaceAll(emailTemplate, "${merchantName}", merchantName)
	emailTemplate = strings.ReplaceAll(emailTemplate, "${verificationLink}", verificationLink)

	return emailTemplate
}
