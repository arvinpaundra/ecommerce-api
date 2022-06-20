package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Category struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null;unique" json:"category_name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Category) Prepare() {
	c.ID = 0
	c.Name = html.EscapeString(strings.TrimSpace(c.Name))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return errors.New("Category name is required")
	}

	return nil
}

func (c *Category) SaveCategory(db *gorm.DB) (*Category, error) {
	var err error
	err = db.Debug().Create(&c).Error

	if err != nil {
		return &Category{}, err
	}

	return c, nil
}

func (c *Category) FindAllCategories(db *gorm.DB) (*[]Category, error) {
	var err error
	categories := []Category{}
	err = db.Debug().Select("name, id").Limit(100).Find(&categories).Error

	if err != nil {
		return &[]Category{}, errors.New("Got 0 result.")
	}

	return &categories, nil
}

func (c *Category) FindCategoryByID(db *gorm.DB, categoryId uint32) (*Category, error) {
	var err error
	err = db.Debug().Model(&Category{}).Where("id = ?", categoryId).Take(&c).Error

	if err != nil {
		return &Category{}, errors.New("Product you have search not exist.")
	}

	return c, nil
}

func (c *Category) UpdateCategory(db *gorm.DB) (*Category, error) {
	var err error
	err = db.Debug().Model(&Category{}).Where("id = ?", c.ID).Updates(Category{Name: c.Name, UpdatedAt: time.Now()}).Error

	if err != nil {
		return &Category{}, err
	}

	return c, nil
}

func (c *Category) DeleteCategory(db *gorm.DB, categoryId uint32) (int64, error) {
	db = db.Debug().Model(&Category{}).Where("id = ?", categoryId).Take(&Category{}).Delete(&Category{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Category not found")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}
