package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	ID         uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"product_name"`
	Stock      uint16    `json:"stock"`
	Price      uint64    `json:"price"`
	Desc       string    `json:"desc"`
	Category   Category  `json:"category"`
	CategoryID uint32    `gorm:"not null" json:"category_id"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Product) Prepare() {
	p.Name = html.EscapeString(strings.TrimSpace(p.Name))
	p.Category = Category{}
	p.Desc = html.EscapeString(strings.TrimSpace(p.Desc))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Product) Validate() error {
	if p.Name == "" {
		return errors.New("Product name is required")
	}

	if p.Stock < 1 {
		return errors.New("Product stock should be available")
	}

	if p.Price < 1 {
		return errors.New("Product price invalid")
	}

	if p.CategoryID < 1 {
		return errors.New("Product should have category")
	}

	return nil
}

func (p *Product) SaveProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Create(&p).Error

	if err != nil {
		return &Product{}, err
	}

	if p.CategoryID > 0 {
		err = db.Debug().Model(&Category{}).Where("id = ?", p.CategoryID).Take(&p.Category).Error

		if err != nil {
			return &Product{}, err
		}
	}

	return p, nil
}

func (p *Product) FindAllProducts(db *gorm.DB) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Table("products").Limit(100).Find(&products).Error

	if err != nil {
		return &[]Product{}, err
	}

	if len(products) < 1 {
		return &[]Product{}, errors.New("Got 0 result.")
	}

	if len(products) > 0 {
		for i := range products {
			err = db.Debug().Model(&Category{}).Where("id = ?", products[i].CategoryID).Take(&products[i].Category).Error
		}

		if err != nil {
			return &[]Product{}, err
		}
	}

	return &products, nil
}

func (p *Product) FindProductByID(db *gorm.DB, productId uint64) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Select("products.id, products.name, products.price, products.stock, products.desc, products.category_id, products.created_at, products.Updated_at").Where("products.id = ?", productId).Take(&p).Error

	if err != nil {
		return &Product{}, errors.New("Product you have search not exist.")
	}

	if p.ID != 0 {
		err = db.Debug().Model(&Category{}).Where("id = ?", p.CategoryID).Take(&p.Category).Error
	}

	if err != nil {
		return &Product{}, errors.New("Product you have search not exist.")
	}

	return p, nil
}

func (p *Product) FindProductByCategory(db *gorm.DB, category uint32) (*[]Product, error) {
	var err error
	products := []Product{}
	err = db.Debug().Model(&Product{}).Where("category_id = ?", category).Find(&products).Error

	if err != nil {
		return &[]Product{}, err
	}

	if err != nil {
		return &[]Product{}, err
	}

	if len(products) > 0 {
		for i := range products {
			err = db.Debug().Model(&Category{}).Where("id = ?", products[i].CategoryID).Take(&products[i].Category).Error
		}

		if err != nil {
			return &[]Product{}, err
		}
	}

	return &products, nil
}

func (p *Product) UpdateProduct(db *gorm.DB) (*Product, error) {
	var err error
	err = db.Debug().Model(&Product{}).Where("id = ?", p.ID).Updates(Product{Name: p.Name, Stock: p.Stock, Price: p.Price, CategoryID: p.CategoryID, Desc: p.Desc, UpdatedAt: time.Now()}).Error

	if err != nil {
		return &Product{}, err
	}

	if p.ID > 0 {
		err = db.Debug().Model(&Category{}).Select("name").Where("id = ?", p.CategoryID).Take(&p.Category).Error

		if err != nil {
			return &Product{}, err
		}
	}

	return p, nil
}

func (p *Product) DeleteProduct(db *gorm.DB, productID uint64) (int64, error) {
	db = db.Debug().Model(&Product{}).Where("id = ?", productID).Take(&Product{}).Delete(&Product{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Product not found")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}

func (p *Product) BeforeUpdateProudct(db *gorm.DB, productId uint64, paramStock uint16) (err error) {
	db.Debug().Model(&Product{}).Where("id = ?", productId).Take(&p)

	if p.Stock < paramStock {
		return errors.New("Stock is not sufficient")
	}

	return nil
}

func (p *Product) AfterUpdateProduct(db *gorm.DB, productId uint64, stockOut uint16) (err error) {
	err = db.Model(&Product{}).Where("id = ?", productId).Take(&p).Update(&Product{Stock: p.Stock - stockOut}).Where("id = ?", productId).Take(&p).Error

	if err != nil {
		return err
	}

	return nil
}
