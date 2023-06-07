package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Inventory is used by pop to map your inventories database table to your go code.
type Inventory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	BookID    string    `json:"book_id" db:"book_id"`
	Qty       int       `json:"qty" db:"qty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Book      *Book     `belongs_to:"books"`
}

// String is not required by pop and may be deleted
func (i Inventory) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Inventories is not required by pop and may be deleted
type Inventories []Inventory

// String is not required by pop and may be deleted
func (i Inventories) String() string {
	ji, _ := json.Marshal(i)
	return string(ji)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (i *Inventory) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsPresent{Field: i.Qty, Name: "Qty"},
		&validators.StringIsPresent{Field: i.BookID, Name: "BookID"},
		// &validators.IntIsPresent{Field: b.Status, Name: "Status"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (i *Inventory) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (i *Inventory) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
