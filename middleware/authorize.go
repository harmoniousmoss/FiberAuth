package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"myfibergotemplate/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware verifies the JWT token and checks user permissions
func AuthMiddleware(c *fiber.Ctx) error {
	// Get the JWT from the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Check if the token starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Malformed Authorization header"})
	}

	// Extract the token from the header
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for validation
		secretKey := config.GetEnv("JWT_SECRET", "your_jwt_secret_key")
		return []byte(secretKey), nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT: " + err.Error()})
	}

	if !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT"})
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT claims"})
	}

	// Store user ID and role in the context
	c.Locals("userID", claims["id"])
	c.Locals("userRole", claims["role"])

	return c.Next()
}

// AdminOnlyMiddleware checks if the user is an administrator
func AdminOnlyMiddleware(c *fiber.Ctx) error {
	if c.Locals("userRole") != "administrator" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Administrator only."})
	}
	return c.Next()
}

// OwnDataOrAdminMiddleware allows users to access their own data or administrators to access any data
func OwnDataOrAdminMiddleware(c *fiber.Ctx) error {
	userID := c.Params("id")
	jwtUserID := c.Locals("userID").(string)
	jwtUserRole := c.Locals("userRole").(string)

	// Allow if the user is accessing their own data or if they are an administrator
	if jwtUserID == userID || jwtUserRole == "administrator" {
		return c.Next()
	}

	return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied."})
}
