package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreateProductImage(image *model.ProductImage) error {
	return config.DB.Create(image).Error
}

func BatchCreateProductImages(images []*model.ProductImage) error {
	return config.DB.Create(images).Error
}

func GetProductImageByID(id uint64) (*model.ProductImage, error) {
	var image model.ProductImage
	err := config.DB.Where("id = ?", id).First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func GetProductImageList(page int, pageSize int, isUsed *bool) ([]model.ProductImage, int64, error) {
	var images []model.ProductImage
	var total int64

	query := config.DB.Model(&model.ProductImage{})

	if isUsed != nil {
		query = query.Where("is_used = ?", *isUsed)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

func UpdateProductImage(image *model.ProductImage) error {
	return config.DB.Save(image).Error
}

func DeleteProductImage(id uint64) error {
	return config.DB.Where("id = ?", id).Delete(&model.ProductImage{}).Error
}

func BatchDeleteProductImages(ids []uint64) error {
	return config.DB.Where("id IN (?)", ids).Delete(&model.ProductImage{}).Error
}
