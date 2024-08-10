package handlers

import (
	"context"
	"net/http"
	"time"

	"myfibergotemplate/config"
	"myfibergotemplate/database"
	"myfibergotemplate/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// SignInHandler handles user sign-in and checks their status
func SignInHandler(c *fiber.Ctx) error {
	type SignInRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	var signInReq SignInRequest

	if err := c.BodyParser(&signInReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User

	err := collection.FindOne(ctx, bson.M{"email": signInReq.Email}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signInReq.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect password"})
	}

	// Check if the user's email is verified
	if !user.EmailStatus {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message":        "Email not verified",
			"email_verified": false,
		})
	}

	// Check if the user's account is approved
	if user.Status != models.Approved {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"message":          "Account not approved",
			"account_approved": false,
			"email_verified":   true,
		})
	}

	// Generate JWT token on successful login
	token, err := generateJWTToken(user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":          "Sign-in successful",
		"email_verified":   true,
		"account_approved": true,
		"token":            token,
		"user": fiber.Map{
			"merchant_name":    user.MerchantName,
			"email":            user.Email,
			"role":             user.Role,
			"person_in_charge": user.PersonInCharge,
			"phone_number":     user.PhoneNumber,
			"website":          user.Website,
			"address":          user.Address,
		},
	})
}

// generateJWTToken creates a JWT token for the authenticated user
func generateJWTToken(user models.User) (string, error) {
	secretKey := config.GetEnv("JWT_SECRET", "your_jwt_secret_key")

	// Define the JWT claims
	claims := jwt.MapClaims{
		"id":    user.ID.Hex(),
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
