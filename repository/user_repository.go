package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func GetUserByPhone(phone string) (*model.User, error) {
	var user model.User
	err := config.DB.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByID(id uint64) (*model.User, error) {
	var user model.User
	err := config.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *model.User) error {
	return config.DB.Save(user).Error
}

func GetPendingUsers() ([]model.User, error) {
	var users []model.User
	err := config.DB.Where("auth_status = ?", model.AuthStatusPending).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	err := config.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUsersByPhone(phone string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := config.DB.Model(&model.User{})
	err := query.Where("phone LIKE ?", "%"+phone+"%").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUsersByStatus(status int, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := config.DB.Model(&model.User{})
	err := query.Where("status = ?", status).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUsersWithFilters(phone string, status *int, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := config.DB.Model(&model.User{})

	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}