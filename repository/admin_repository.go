package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func GetAdminByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	err := config.DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func GetAdminByID(id uint64) (*model.Admin, error) {
	var admin model.Admin
	err := config.DB.Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}