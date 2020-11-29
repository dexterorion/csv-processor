package model

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	TimeIntervalHour []time.Time        `bson:"time_interval_hour"`
	CheckinDate      time.Time          `bson:"checkin_date"`
	CheckoutDate     time.Time          `bson:"checkout_date"`
	FareAmount       float64            `bson:"fare_amount"`
	FareName         string             `bson:"fare_name"`
	PaidAmount       float64            `bson:"paid_amount"`
	Discount         float64            `bson:"discount"`
	PaymentData      string             `bson:"payment_data"`
	PaymentMethod    string             `bson:"payment_method"`
	UseType          string             `bson:"use_type"`
	OfferType        string             `bson:"offer_type"`
	Duration         float64            `bson:"duration"`
	IsValid          bool               `bson:"is_valid"`
	ParkingInfo      Parking            `bson:"parking_info"`
	Status           string             `bson:"status"`
	CashRegisterID   int64              `bson:"cash_register_id"`
	Sequence         string             `bson:"sequence"`
	Fiscal           string             `bson:"fiscal"`
	Partial          string             `bson:"partial"`

	Version   int        `bson:"version"`
	Schema    int        `bson:"schema"`
	CreatedAt time.Time  `bson:"created_at"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty"`
}

// SchemaVersion returns the schema version
func (t Transaction) SchemaVersion() int {
	return 1
}

// Validate validates the model
func (t Transaction) Validate() error {
	errs := []string{}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "."))
	}

	return nil
}

type Parking struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
	Slug string `bson:"slug"`
}
