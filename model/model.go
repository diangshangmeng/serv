package model

import (
	"time"

	"voucher-platform/config"

	"github.com/jinzhu/gorm"
)

// AuthStatus constants
const (
	AuthStatusUnverified = 0 // 未验证状态
	AuthStatusPending    = 1 // 等待后台管理员进行验证
	AuthStatusApproved   = 2 // 通过后台管理验证状态
	AuthStatusRejected   = 3 // 认证驳回
)

type User struct {
	gorm.Model
	Phone           string `gorm:"unique;not null"`
	Password        string `gorm:"not null"`
	CityID          uint
	Status          int `gorm:"default:1"`
	AuthStatus      int `gorm:"default:0"`
	AuditRemark     string
	IDCardFront     string
	IDCardBack      string
	BusinessLicense string
	PaymentCode     string
}

type City struct {
	gorm.Model
	Name       string `gorm:"not null"`
	ProvinceID uint   `gorm:"default:1"`
}

type Coupon struct {
	gorm.Model
	UserID        uint    `gorm:"not null"`
	Code          string  `gorm:"unique;not null"`
	OriginalPrice float64 `gorm:"not null"`
	Status        int     `gorm:"default:0"`
}

type Listing struct {
	gorm.Model
	CouponID     uint    `gorm:"not null"`
	ProductID    uint    `gorm:"not null"`
	SellingPrice float64 `gorm:"not null"`
	Status       int     `gorm:"default:0"`
}

type Order struct {
	gorm.Model
	OrderNo             string `gorm:"unique;not null"`
	ProductID           uint   `gorm:"index"`
	ProductTitle        string
	ProductDescription  string
	SellerID            uint
	SellerPhone         string
	ProductImageURL     string `gorm:"size:2048"`
	PaymentCodeImageURL string `gorm:"size:2048"`
	BuyerID             uint   `gorm:"index"`
	BuyerPhone          string
	PaymentVoucher      string `gorm:"size:2048"`
	Status              int    `gorm:"default:0"`
	Price               int64
	ExpiredAt           time.Time
	OrderTime           *time.Time
}

type Admin struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type ProductImage struct {
	gorm.Model
	Path     string
	Name     string
	FileName string
	Tags     string
	IsUsed   bool `gorm:"default:false"`
	AdminID  uint
	Size     int64
}

type PaymentCodeImage struct {
	gorm.Model
	Path        string
	CompanyName string
	FileName    string
	AdminID     uint
	Size        int64
}

type Product struct {
	gorm.Model
	Title           string
	Description     string
	Price           int64
	DisplayImageID  uint
	PaymentImageID  uint
	OwnerID         uint
	OwnerPhone      string
	DisplayImageURL string `gorm:"size:2048"`
	PaymentImageURL string `gorm:"size:2048"`
	Status          int    `gorm:"default:0"`
	PayTime         *time.Time
	LockReason      string
	CityID          uint `gorm:"index"`
	City            City
}

func AutoMigrate() error {
	err := config.DB.AutoMigrate(
		&User{},
		&City{},
		&Coupon{},
		&Listing{},
		&Order{},
		&Admin{},
		&ProductImage{},
		&PaymentCodeImage{},
		&Product{},
	).Error
	if err != nil {
		return err
	}

	err = initCities()
	if err != nil {
		return err
	}

	err = initAdmin()
	if err != nil {
		return err
	}

	return nil
}

func initCities() error {
	var count int
	config.DB.Model(&City{}).Count(&count)
	if count > 0 {
		return nil
	}

	cities := []City{
		{Name: "昆明市", ProvinceID: 1},
		{Name: "大理市", ProvinceID: 1},
		{Name: "曲靖市", ProvinceID: 1},
		{Name: "玉溪市", ProvinceID: 1},
		{Name: "保山市", ProvinceID: 1},
		{Name: "昭通市", ProvinceID: 1},
		{Name: "丽江市", ProvinceID: 1},
		{Name: "普洱市", ProvinceID: 1},
		{Name: "临沧市", ProvinceID: 1},
	}

	for _, city := range cities {
		if err := config.DB.Create(&city).Error; err != nil {
			return err
		}
	}

	return nil
}

func initAdmin() error {
	var count int
	config.DB.Model(&Admin{}).Count(&count)
	if count > 0 {
		return nil
	}

	admin := Admin{
		Username: "admin",
		Password: "$2a$10$kj1dClEJF9gbL6axUER.PeYNajymob4cI9r/I9YxPiiIYsFjGyp4q",
	}

	return config.DB.Create(&admin).Error
}
