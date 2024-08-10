package handlers

import (
	"context"
	"net/http"
	"time"

	"myfibergotemplate/database"
	"myfibergotemplate/libs"
	"myfibergotemplate/models"
	"myfibergotemplate/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// EditUserHandler handles editing user information
func EditUserHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	authUserID := c.Locals("userID").(string)
	authUserRole := c.Locals("userRole").(string)

	if !isAuthorizedToEdit(authUserRole, authUserID, userID) {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	update := createUpdateDocument(updateData, c)

	collection := database.GetMongoClient().Database("talentdevgo").Collection("users")
	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

func isAuthorizedToEdit(role, authUserID, userID string) bool {
	return role == "administrator" || authUserID == userID
}

func createUpdateDocument(updateData models.User, c *fiber.Ctx) bson.M {
	update := bson.M{}

	// Ensure the userEmail is present in the locals before trying to use it
	if userEmail, ok := c.Locals("userEmail").(string); ok && updateData.Email != "" && updateData.Email != userEmail {
		update["email"] = updateData.Email
		update["email_status"] = false
		verificationToken := utils.GenerateVerificationToken()
		update["verification_token"] = verificationToken
		sendVerificationEmail(updateData.MerchantName, updateData.Email, verificationToken)
	}

	// Handle other fields
	if updateData.MerchantName != "" {
		update["merchant_name"] = updateData.MerchantName
	}
	if updateData.PersonInCharge != "" {
		update["person_in_charge"] = updateData.PersonInCharge
	}
	if updateData.PhoneNumber != "" {
		update["phone_number"] = updateData.PhoneNumber
	}
	if updateData.Website != "" {
		update["website"] = updateData.Website
	}
	if updateData.Address != "" {
		update["address"] = updateData.Address
	}
	if updateData.TermsAndConditions {
		update["terms_and_conditions"] = updateData.TermsAndConditions
	}
	if updateData.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
		update["password"] = string(hashedPassword)
	}

	update["updated_at"] = time.Now()
	return update
}

func sendVerificationEmail(merchantName, email, token string) {
	verificationLink := "http://localhost:8080/api/verify-email?token=" + token
	emailBody := `<p>Dear ` + merchantName + `,</p>
				  <p>You have requested to change your email address. Please verify your new email by clicking the link below:</p>
				  <p><a href="` + verificationLink + `">Verify Email Address</a></p>
				  <p>Kind regards,<br>The TalentDev Team</p>`
	libs.SendEmail([]string{email}, "Email Verification - TalentDev ID", emailBody)
}
