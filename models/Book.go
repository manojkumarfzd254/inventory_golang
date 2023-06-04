package models

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Book is used by pop to map your books database table to your go code.
type Book struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	CategoryID  string       `json:"category_id" db:"category_id"`
	Title       string       `json:"title" db:"title"`
	BookNo      string       `json:"book_no" db:"book_no"`
	Author      string       `json:"author" db:"author"`
	Picture     binding.File `db:"-" form:"picture"`
	PicturePath string       `json:"picture_path" db:"picture_path"`
	Price       float64      `json:"price" db:"price"`
	Status      int          `json:"status" db:"status"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	Category    *Category    `belongs_to:"categories"`
}

// String is not required by pop and may be deleted
func (b Book) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Books is not required by pop and may be deleted
type Books []Book

// String is not required by pop and may be deleted
func (b Books) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

func (b *Book) Create(tx *pop.Connection) (*validate.Errors, error) {

	dir := filepath.Join(".", "uploads/books")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}

	f, err := os.Create(filepath.Join(dir, b.Picture.Filename))
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	defer f.Close()
	_, err = io.Copy(f, b.Picture)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	b.PicturePath = "/" + filepath.Join(dir, b.Picture.Filename)
	return tx.ValidateAndCreate(b)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *Book) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *Book) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *Book) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
