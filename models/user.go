package models

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is a generated model from buffalo-auth, it serves as the base for username/password authentication.
type User struct {
	ID                   uuid.UUID    `json:"id" db:"id"`
	CreatedAt            time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at" db:"updated_at"`
	Email                string       `json:"email" db:"email"`
	PasswordHash         string       `json:"password_hash" db:"password_hash"`
	Name                 string       `json:"name" db:"name"`
	Mobile               string       `json:"mobile" db:"mobile"`
	Address              string       `json:"address" db:"address"`
	Password             string       `json:"-" db:"-"`
	PasswordConfirmation string       `json:"-" db:"-"`
	Profile              binding.File `db:"-" form:"profile"`
	ProfilePath          string       `json:"profile_path" db:"profile_path"`
}

// Create wraps up the pattern of encrypting the password and
// running validations. Useful when writing tests.
func (u *User) Create(tx *pop.Connection) (*validate.Errors, error) {
	u.Email = strings.ToLower(u.Email)
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndCreate(u)
}

func (u *User) Update(tx *pop.Connection) (*validate.Errors, error) {
	// if !u.Profile.Valid() {

	// }
	if u.Profile.Valid() {
		dir := filepath.Join(".", "uploads")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		currentTime := time.Now()
		// Convert the current time to a numerical representation
		timeNumber := currentTime.Unix()
		ext := filepath.Ext(u.Profile.Filename)
		filename := strconv.FormatInt(timeNumber, 10) + ext
		filePath := filepath.Join(dir, filename)

		f, err := os.Create(filepath.Join(dir, u.Profile.Filename))
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		defer f.Close()
		_, err = io.Copy(f, u.Profile)
		if err != nil {
			return validate.NewErrors(), errors.WithStack(err)
		}
		if u.ProfilePath != "" && u.ProfilePath != "/"+filePath {
			ProfilePath := strings.TrimLeft(u.ProfilePath, "/")
			if err := os.Remove(ProfilePath); err != nil {
				return validate.NewErrors(), errors.WithStack(err)
			}
		}
		u.ProfilePath = "/" + filepath.Join(dir, u.Profile.Filename)
	}
	ph, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return validate.NewErrors(), errors.WithStack(err)
	}
	u.PasswordHash = string(ph)
	return tx.ValidateAndUpdate(u)
}

// String is not required by pop and may be deleted
func (u User) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	ju, _ := json.Marshal(u)
	return string(ju)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.PasswordHash, Name: "PasswordHash"},
		// check to see if the email address is already taken:
		&validators.FuncValidator{
			Field:   u.Email,
			Name:    "Email",
			Message: "%s is already taken",
			Fn: func() bool {
				var b bool
				q := tx.Where("email = ?", u.Email)
				if u.ID != uuid.Nil {
					q = q.Where("id != ?", u.ID)
				}
				b, err = q.Exists(u)
				if err != nil {
					return false
				}
				return !b
			},
		},
	), err
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (u *User) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	var err error
	return validate.Validate(
		&validators.StringIsPresent{Field: u.Password, Name: "Password"},
		&validators.StringsMatch{Name: "Password", Field: u.Password, Field2: u.PasswordConfirmation, Message: "Password does not match confirmation"},
	), err
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
