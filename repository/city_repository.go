package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func GetAllCities() ([]model.City, error) {
	var cities []model.City
	err := config.DB.Find(&cities).Error
	if err != nil {
		return nil, err
	}
	return cities, nil
}

func GetCityByID(id uint64) (*model.City, error) {
	var city model.City
	err := config.DB.Where("id = ?", id).First(&city).Error
	if err != nil {
		return nil, err
	}
	return &city, nil
}