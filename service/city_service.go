package service

import (
	"voucher-platform/repository"
)

func GetCityList() ([]interface{}, error) {
	cities, err := repository.GetAllCities()
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(cities))
	for i, city := range cities {
		result[i] = city
	}

	return result, nil
}
