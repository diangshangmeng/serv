package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreateProduct(product *model.Product) error {
	return config.DB.Create(product).Error
}

func GetProductByID(id uint64) (*model.Product, error) {
	var product model.Product
	err := config.DB.Preload("City").Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetProductList(page int, pageSize int, cityID *uint) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := config.DB.Model(&model.Product{}).Where("status != ?", 0)

	if cityID != nil {
		query = query.Where("city_id = ?", *cityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("City").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func GetProductListForDashboard(page int, pageSize int, cityID *uint) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := config.DB.Model(&model.Product{})

	if cityID != nil {
		query = query.Where("city_id = ?", *cityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("City").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func UpdateProduct(product *model.Product) error {
	return config.DB.Save(product).Error
}

func DeleteProduct(id uint64) error {
	return config.DB.Where("id = ?", id).Delete(&model.Product{}).Error
}

func GetProductsByOwnerPhone(ownerPhone string, page int, pageSize int, cityID *uint) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	query := config.DB.Model(&model.Product{}).Where("owner_phone = ?", ownerPhone)

	if cityID != nil {
		query = query.Where("city_id = ?", *cityID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("City").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
