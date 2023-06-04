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
	"github.com/gobuffalo/validate/v3/validators"
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
	Price       string       `json:"price" db:"price"`
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
	// Create the directory for storing book pictures

	// uuid := uuid.New()
	// Generate a unique filename for the picture
	if b.Picture.Valid() {
		dir := filepath.Join(".", "uploads", "books")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}

		filename := uuid.NamespaceDNS.String() + filepath.Ext(b.Picture.Filename)
		filePath := filepath.Join(dir, filename)

		// Create the file and copy the picture data
		file, err := os.Create(filePath)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		defer file.Close()

		_, err = io.Copy(file, b.Picture)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}

		// Set the PicturePath field of the Book struct
		b.PicturePath = "/" + filePath
	}
	// Validate and create the book record in the database
	verrs, err := tx.ValidateAndCreate(b)
	if err != nil {
		return verrs, errors.WithStack(err)
	}

	return verrs, nil
}

func (b *Book) Update(tx *pop.Connection) (*validate.Errors, error) {
	// Create the directory for storing book pictures

	// uuid := uuid.New()
	// Generate a unique filename for the picture
	if b.Picture.Valid() {
		dir := filepath.Join(".", "uploads", "books")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}

		filename := uuid.NamespaceDNS.String() + filepath.Ext(b.Picture.Filename)
		filePath := filepath.Join(dir, filename)

		// Create the file and copy the picture data
		file, err := os.Create(filePath)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		defer file.Close()

		_, err = io.Copy(file, b.Picture)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}

		// Set the PicturePath field of the Book struct
		b.PicturePath = "/" + filePath
	}
	// Validate and create the book record in the database
	verrs, err := tx.ValidateAndUpdate(b)
	if err != nil {
		return verrs, errors.WithStack(err)
	}

	return verrs, nil
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *Book) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.Title, Name: "Title"},
		&validators.StringIsPresent{Field: b.CategoryID, Name: "CategoryID"},
		&validators.StringIsPresent{Field: b.BookNo, Name: "BookNo"},
		&validators.StringIsPresent{Field: b.Author, Name: "Author"},
		&validators.StringIsPresent{Field: b.Price, Name: "Price"},
		// &validators.IntIsPresent{Field: b.Status, Name: "Status"},
	), nil
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
