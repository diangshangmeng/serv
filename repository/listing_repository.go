package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreateListing(listing *model.Listing) error {
	return config.DB.Create(listing).Error
}

func GetListingByID(id uint64) (*model.Listing, error) {
	var listing model.Listing
	err := config.DB.Where("id = ?", id).First(&listing).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

func GetListingByProductID(productID uint64) (*model.Listing, error) {
	var listing model.Listing
	err := config.DB.Where("product_id = ?", productID).First(&listing).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

func GetListingsByCity(cityID uint) ([]model.Listing, error) {
	var listings []model.Listing
	err := config.DB.Joins("JOIN coupons ON listings.coupon_id = coupons.id").
		Joins("JOIN users ON coupons.user_id = users.id").
		Where("users.city_id = ? AND listings.status = ?", cityID, 0).
		Find(&listings).Error
	if err != nil {
		return nil, err
	}
	return listings, nil
}

func GetListingsByUserID(userID uint64) ([]model.Listing, error) {
	var listings []model.Listing
	err := config.DB.Joins("JOIN coupons ON listings.coupon_id = coupons.id").
		Where("coupons.user_id = ?", userID).
		Find(&listings).Error
	if err != nil {
		return nil, err
	}
	return listings, nil
}

func UpdateListing(listing *model.Listing) error {
	return config.DB.Save(listing).Error
}

func UpdateListingStatus(listingID uint64, status int) error {
	return config.DB.Model(&model.Listing{}).Where("id = ?", listingID).Update("status", status).Error
}