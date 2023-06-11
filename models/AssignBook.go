package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// AssignBook is used by pop to map your assign_books database table to your go code.
type AssignBook struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CustomerID string    `json:"customer_id" db:"customer_id"`
	BookID     string    `json:"book_id" db:"book_id"`
	AssignDate string    `json:"assign_date" db:"assign_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a AssignBook) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// AssignBooks is not required by pop and may be deleted
type AssignBooks []AssignBook

// String is not required by pop and may be deleted
func (a AssignBooks) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *AssignBook) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.CustomerID, Name: "CustomerID"},
		&validators.StringIsPresent{Field: a.BookID, Name: "BookID"},
		&validators.StringIsPresent{Field: a.AssignDate, Name: "AssignDate"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *AssignBook) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *AssignBook) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
