package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Payment struct {
	ID        uint8     `gorm:"primary_key;auto_increment" json:"id"`
	BankName  string    `gorm:"size:100;not null;unique" json:"bank_name"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Payment) Prepare() {
	p.ID = 0
	p.BankName = html.EscapeString(strings.TrimSpace(p.BankName))
}

func (p *Payment) Validate() error {
	if p.BankName == "" {
		return errors.New("Name of the Bank is required")
	}

	return nil
}

func (p *Payment) SavePayment(db *gorm.DB) (*Payment, error) {
	var err error
	err = db.Debug().Model(&Payment{}).Create(&p).Error

	if err != nil {
		return &Payment{}, err
	}

	return p, nil
}

func (p *Payment) FindAllPaymentsMethod(db *gorm.DB) (*[]Payment, error) {
	var err error
	payments := []Payment{}
	err = db.Debug().Model(&Payment{}).Limit(100).Find(&payments).Error

	if err != nil {
		return &[]Payment{}, err
	}

	return &payments, nil
}

func (p *Payment) UpdatePayment(db *gorm.DB) (*Payment, error) {
	var err error
	err = db.Debug().Model(&Payment{}).Where("id = ?", p.ID).Update(Payment{BankName: p.BankName}).Error

	if err != nil {
		return &Payment{}, err
	}

	return p, nil
}

func (p *Payment) DeletePayment(db *gorm.DB, paymentId uint8) (int64, error) {
	db = db.Debug().Model(&Payment{}).Where("id = ?", paymentId).Take(&Payment{}).Delete(&Payment{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Payment method not found")
		}

		return 0, db.Error
	}

	return db.RowsAffected, nil
}
