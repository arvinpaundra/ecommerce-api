package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Checkout struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	CustomerID uint32    `json:"customer_id"`
	PaymentID  uint8     `json:"payment_id"`
	CartID     uint64    `json:"cart_id"`
	Customer   Customer  `json:"customer_details"`
	Payment    Payment   `json:"payment"`
	Cart       Cart      `json:"cart"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Checkout) Prepare() {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Checkout) Validate() error {
	if c.CartID < 1 {
		return errors.New("Cart details not found")
	}

	if c.CustomerID < 1 {
		return errors.New("Customer details not found")
	}

	if c.PaymentID < 1 {
		return errors.New("Paymnet details not found")
	}

	return nil
}

func (c *Checkout) CreateCustomerCheckout(db *gorm.DB) (*Checkout, error) {
	var err error
	err = db.Debug().Model(&Checkout{}).Create(&c).Error

	if err != nil {
		return &Checkout{}, err
	}

	if c.CustomerID > 0 {
		db := db.Debug().Model(&Customer{}).Where("id = ?", c.CustomerID).Take(&c.Customer)

		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return &Checkout{}, errors.New("Customer not found")
			}
			return &Checkout{}, db.Error
		}
	}

	if c.PaymentID > 0 {
		db := db.Debug().Model(&Payment{}).Where("id = ?", c.PaymentID).Take(&c.Payment)

		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return &Checkout{}, errors.New("Payment not found")
			}
			return &Checkout{}, db.Error
		}
	}

	if c.CartID > 0 {
		db := db.Debug().Model(&Cart{}).Where("id = ?", c.CartID).Take(&c.Cart)

		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return &Checkout{}, errors.New("Cart not found")
			}
			return &Checkout{}, db.Error
		}
	}

	if c.Cart.ProductID > 0 {
		db := db.Debug().Model(&Product{}).Where("id = ?", c.Cart.ProductID).Take(&c.Cart.Product)

		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return &Checkout{}, errors.New("Product not found")
			}
			return &Checkout{}, db.Error
		}
	}

	if c.Cart.Product.CategoryID > 0 {
		db := db.Debug().Model(&Category{}).Where("id = ?", c.Cart.Product.CategoryID).Take(&c.Cart.Product.Category)

		if db.Error != nil {
			if gorm.IsRecordNotFoundError(db.Error) {
				return &Checkout{}, errors.New("Category not found")
			}
			return &Checkout{}, db.Error
		}
	}

	return c, nil
}

func (c *Checkout) FindCustomerCheckout(db *gorm.DB, customerId uint64) (*[]Checkout, error) {
	var err error
	checkouts := []Checkout{}
	err = db.Debug().Model(&Checkout{}).Where("customer_id = ?", customerId).Limit(100).Find(&checkouts).Error

	if err != nil {
		return &[]Checkout{}, err
	}

	if len(checkouts) > 0 {
		for i := range checkouts {
			err = db.Debug().Model(&Customer{}).Where("id = ?", checkouts[i].CustomerID).Take(&checkouts[i].Customer).Error
		}

		if err != nil {
			return &[]Checkout{}, errors.New("Not found customer")
		}

		for i := range checkouts {
			err = db.Debug().Model(&Payment{}).Where("id = ?", checkouts[i].PaymentID).Take(&checkouts[i].Payment).Error
		}

		if err != nil {
			return &[]Checkout{}, errors.New("Not found payment")
		}

		for i := range checkouts {
			err = db.Debug().Model(&Cart{}).Where("id = ?", checkouts[i].CartID).Take(&checkouts[i].Cart).Error
		}

		if err != nil {
			return &[]Checkout{}, errors.New("Not found cart")
		}

		for i := range checkouts {
			err = db.Debug().Model(&Product{}).Where("id = ?", checkouts[i].Cart.ProductID).Take(&checkouts[i].Cart.Product).Error
		}

		if err != nil {
			return &[]Checkout{}, errors.New("Not found product")
		}

		for i := range checkouts {
			err = db.Debug().Model(&Category{}).Where("id = ?", checkouts[i].Cart.Product.CategoryID).Take(&checkouts[i].Cart.Product.Category).Error
		}

		if err != nil {
			return &[]Checkout{}, errors.New("Not found category")
		}
	}

	if len(checkouts) < 1 {
		return &[]Checkout{}, errors.New("You have 0 checkout item.")
	}

	return &checkouts, nil
}
