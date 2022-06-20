package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Cart struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	CustomerID uint32    `json:"customer_id"`
	ProductID  uint64    `json:"product_id"`
	Product    Product   `json:"product"`
	Qty        uint16    `json:"qty"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Cart) Prepare() {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Cart) Validate() error {
	if c.CustomerID < 1 {
		return errors.New("Customer have to be provided")
	}

	if c.ProductID < 1 {
		return errors.New("Product item is required")
	}

	if c.Qty < 1 {
		return errors.New("At least one amount of product have to added into cart")
	}

	return nil
}

func (c *Cart) AddToCart(db *gorm.DB) (*Cart, error) {
	var err error
	err = db.Debug().Model(&Cart{}).Create(&c).Error

	if err != nil {
		return &Cart{}, err
	}

	if c.CustomerID < 0 {
		return &Cart{}, errors.New("Customer not provided")
	}

	if c.ProductID > 0 {
		err = db.Debug().Model(&Product{}).Where("id = ?", c.ProductID).Take(&c.Product).Error

		if err != nil {
			return &Cart{}, err
		}

		err = db.Debug().Model(&Category{}).Where("id = ?", c.Product.CategoryID).Take(&c.Product.Category).Error

		if err != nil {
			return &Cart{}, err
		}
	}

	return c, nil
}

func (c *Cart) DeleteCart(db *gorm.DB, cartId uint64, customerId uint32) (int64, error) {
	db = db.Debug().Model(&Cart{}).Where("id = ? AND customer_id = ?", cartId, customerId).Take(&Cart{}).Delete(&Cart{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Cart not found")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}

// func (c *Cart) FindCustomerCart(db *gorm.DB, customerId uint32) (*[]Cart, error) {
// 	var err error
// 	carts := []Cart{}
// 	err = db.Debug().Model(&Cart{}).Where("customer_id = ?", customerId).Limit(100).Find(&carts).Error

// 	if err != nil {
// 		return &[]Cart{}, err
// 	}

// 	if len(carts) > 0 {
// 		for i := range carts {
// 			err = db.Debug().Model(&Product{}).Where("id = ?", carts[i].ProductID).Limit(100).Take(&carts[i].Product).Error

// 			if err != nil {
// 				return &[]Cart{}, err
// 			}
// 		}
// 	} else {
// 		return &[]Cart{}, errors.New("You have 0 item in your cart")
// 	}

// 	return &carts, nil
// }

func (c *Cart) FindCustomerCart(db *gorm.DB, customerId uint32) (*[]Cart, error) {
	var err error
	carts := []Cart{}
	err = db.Debug().Model(&Cart{}).Where("customer_id = ?", customerId).Find(&carts).Error

	if err != nil {
		return &[]Cart{}, err
	}

	if len(carts) > 0 {
		for i := range carts {
			err = db.Debug().Model(&Product{}).Where("id = ?", carts[i].ProductID).Find(&carts[i].Product).Error
		}

		if err != nil {
			return &[]Cart{}, err
		}

		for i := range carts {
			err = db.Debug().Model(&Category{}).Where("id = ?", carts[i].Product.CategoryID).Find(&carts[i].Product.Category).Error
		}

		if err != nil {
			return &[]Cart{}, err
		}
	}

	if len(carts) < 1 {
		return &[]Cart{}, errors.New("You have 0 item in your cart.")
	}

	return &carts, nil
}
