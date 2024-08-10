package handlers

import (
	"context"
	"crypto/rand"
	"math/big"
	"net/http"
	"time"

	"myfibergotemplate/database"
	"myfibergotemplate/libs"
	"myfibergotemplate/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var failedAttemptsCache = make(map[string]int)       // Cache to track failed password reset attempts
var lockedAccountsCache = make(map[string]time.Time) // Cache to track locked accounts

const maxFailedAttempts = 5              // Maximum allowed failed attempts
const lockoutDuration = 15 * time.Minute // Duration of the lockout period

// ForgotPasswordHandler handles password reset requests
func ForgotPasswordHandler(c *fiber.Ctx) error {
	type ForgotPasswordRequest struct {
		Email string `json:"email" validate:"required,email"`
	}

	var req ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	email := req.Email

	// Rate limiting for password reset requests
	if lockoutTime, exists := lockedAccountsCache[email]; exists {
		if time.Until(lockoutTime) > 0 {
			remainingTime := time.Until(lockoutTime).Minutes()
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many failed attempts, please try again later",
				"retry_after": int(remainingTime),
			})
		} else {
			delete(lockedAccountsCache, email)
			failedAttemptsCache[email] = 0 // Reset failed attempts counter
		}
	}

	// Get a reference to the "users" collection in the MongoDB database
	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	// Find the user by email
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		// Increment the failed attempts counter
		failedAttemptsCache[email]++
		if failedAttemptsCache[email] >= maxFailedAttempts {
			lockedAccountsCache[email] = time.Now().Add(lockoutDuration)
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many failed attempts, please try again later",
				"retry_after": int(lockoutDuration.Minutes()),
			})
		}
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Check if the user's email is verified and status is approved
	if !user.EmailStatus || user.Status != models.Approved {
		// Increment the failed attempts counter
		failedAttemptsCache[email]++
		if failedAttemptsCache[email] >= maxFailedAttempts {
			lockedAccountsCache[email] = time.Now().Add(lockoutDuration)
			return c.Status(http.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many failed attempts, please try again later",
				"retry_after": int(lockoutDuration.Minutes()),
			})
		}
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "User is not approved or email is not verified"})
	}

	// Generate a new random password with numbers, letters, and special characters
	newPassword := generateSecurePassword(8)

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate password"})
	}

	// Update the user's password in the database
	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}
	_, err = collection.UpdateOne(ctx, bson.M{"email": email}, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	// Compose the email subject and body
	emailSubject := "Forgot password - TalentDev ID"
	emailBody := `<p>Dear ` + user.MerchantName + `,</p>
				  <p>You requested a new password on TalentDev!</p>
				  <p>Here is your new password: ` + newPassword + `</p>
				  <p>Kind regards,<br>The TalentDev Team</p>`

	// Send the new password to the user's email
	err = libs.SendEmail([]string{user.Email}, emailSubject, emailBody)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
	}

	// Return a success message
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "A new password has been sent to your email"})
}

// generateSecurePassword generates a random secure password of a given length
func generateSecurePassword(length int) string {
	const (
		upperLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowerLetters = "abcdefghijklmnopqrstuvwxyz"
		numbers      = "0123456789"
		specialChars = "!@#$%^&*()-_+=<>?{}[]|"
	)

	allChars := upperLetters + lowerLetters + numbers + specialChars
	password := make([]byte, length)

	// Ensure at least one character from each category
	password[0] = upperLetters[randomInt(len(upperLetters))]
	password[1] = lowerLetters[randomInt(len(lowerLetters))]
	password[2] = numbers[randomInt(len(numbers))]
	password[3] = specialChars[randomInt(len(specialChars))]

	// Fill the rest of the password with random characters from all categories
	for i := 4; i < length; i++ {
		password[i] = allChars[randomInt(len(allChars))]
	}

	// Shuffle the characters to ensure randomness
	shuffle(password)

	return string(password)
}

// randomInt generates a random integer in the range [0, max)
func randomInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

// shuffle shuffles a slice of bytes
func shuffle(password []byte) {
	for i := len(password) - 1; i > 0; i-- {
		j := randomInt(i + 1)
		password[i], password[j] = password[j], password[i]
	}
}
