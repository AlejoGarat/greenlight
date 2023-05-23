package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"greenlight/pkg/validator"

	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = User{}

type User struct {
	ID        int64     `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  Password  `json:"-" db:"password_hash"`
	Activated bool      `json:"activated" db:"activated"`
	Version   int       `json:"-" db:"version"`
}

type Password struct {
	Plaintext *string `json:"plain_text"`
	Hash      []byte  `json:"hash"`
}

func (u User) IsAnonymous() bool {
	return u.Email == ""
}

func (p Password) Value() (driver.Value, error) {
	return p.Hash, nil
}

func (p *Password) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		p.Hash = v
	case string:
		p.Hash = []byte(v)
	default:
		return errors.New("incompatible type for Password")
	}
	return nil
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.Plaintext = &plaintextPassword
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
