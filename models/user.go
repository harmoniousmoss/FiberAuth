package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	Administrator Role = "administrator"
	Merchant      Role = "merchant"
)

type Status string

const (
	Approved Status = "approved"
	Pending  Status = "pending"
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id"`
	MerchantName       string             `json:"merchant_name" bson:"merchant_name" validate:"required"`
	Status             Status             `json:"status" bson:"status" validate:"omitempty,oneof=approved pending"`
	Email              string             `json:"email" bson:"email" validate:"required,email"`
	EmailStatus        bool               `json:"email_status" bson:"email_status"`
	Role               Role               `json:"role" bson:"role" validate:"required,oneof=administrator merchant"`
	PersonInCharge     string             `json:"person_in_charge" bson:"person_in_charge" validate:"required"`
	PhoneNumber        string             `json:"phone_number" bson:"phone_number"`
	Website            string             `json:"website" bson:"website"`
	Address            string             `json:"address" bson:"address"`
	Password           string             `json:"password,omitempty" bson:"password"`
	VerificationToken  string             `json:"-" bson:"verification_token,omitempty"`
	TermsAndConditions bool               `json:"terms_and_conditions" bson:"terms_and_conditions" validate:"required"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}
