package service

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
	"github.com/jinzhu/gorm"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func SendCode(phone string) error {
	code := util.GenerateCode()
	return util.SendSMS(phone, code)
}

func Register(phone, password string) (*model.User, error) {
	_, err := repository.GetUserByPhone(phone)
	if err == nil {
		return nil, util.NewBizError(util.ErrCodeUserExists, "用户已存在")
	}
	if err != nil && err != gorm.ErrRecordNotFound && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Phone:    phone,
		Password: string(hashedPassword),
	}

	err = repository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Login(phone, password string) (string, *model.User, error) {
	user, err := repository.GetUserByPhone(phone)
	if err != nil {
		if err == gorm.ErrRecordNotFound || strings.Contains(err.Error(), "record not found") {
			return "", nil, util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
		}
		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, util.NewBizError(util.ErrCodePasswordError, "密码错误")
	}

	token, err := util.GenerateToken(uint64(user.ID), user.Phone, user.CityID)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}