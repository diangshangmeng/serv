package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreatePaymentCodeImage(image *model.PaymentCodeImage) error {
	return config.DB.Create(image).Error
}

func GetPaymentCodeImageByID(id uint64) (*model.PaymentCodeImage, error) {
	var image model.PaymentCodeImage
	err := config.DB.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func GetPaymentCodeImageList(page int, pageSize int) ([]model.PaymentCodeImage, int64, error) {
	var images []model.PaymentCodeImage
	var total int64

	query := config.DB.Model(&model.PaymentCodeImage{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

func UpdatePaymentCodeImage(image *model.PaymentCodeImage) error {
	return config.DB.Save(image).Error
}

func DeletePaymentCodeImage(id uint64) error {
	return config.DB.Where("id = ?", id).Delete(&model.PaymentCodeImage{}).Error
}
