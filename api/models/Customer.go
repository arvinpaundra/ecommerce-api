package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Fullname  string    `gorm:"size:255;not null" json:"fullname"`
	Address   string    `gorm:"size:255;not null" json:"address"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (c *Customer) BeforeSave() error {
	hashedPassword, err := Hash(c.Password)

	if err != nil {
		return err
	}

	c.Password = string(hashedPassword)
	return nil
}

func (c *Customer) Prepare() {
	c.ID = 0
	c.Fullname = html.EscapeString(strings.TrimSpace(c.Fullname))
	c.Address = html.EscapeString(strings.TrimSpace(c.Address))
	c.Email = html.EscapeString(strings.TrimSpace(c.Email))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Customer) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if c.Email == "" {
			return errors.New("Email field is required")
		}

		if c.Password == "" {
			return errors.New("Password field is required")
		}

		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil

	default:
		if c.Fullname == "" {
			return errors.New("Enter your fullname")
		}

		if c.Address == "" {
			return errors.New("Enter your address")
		}

		if c.Email == "" {
			return errors.New("Enter your email")
		}

		if c.Password == "" {
			return errors.New("Enter your password")
		}

		if err := checkmail.ValidateFormat(c.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	}
}

func (c *Customer) SaveCustomer(db *gorm.DB) (*Customer, error) {
	var err error
	err = db.Debug().Model(&Customer{}).Create(&c).Error

	if err != nil {
		return &Customer{}, err
	}

	return c, nil
}
