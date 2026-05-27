package service

import (
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func SubmitAuth(userID uint64, cityID uint, idCardFront, idCardBack, businessLicense, paymentCode string) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.CityID = cityID
	user.IDCardFront = idCardFront
	user.IDCardBack = idCardBack
	user.BusinessLicense = businessLicense
	user.PaymentCode = paymentCode
	user.AuthStatus = model.AuthStatusPending

	return repository.UpdateUser(user)
}

func GetAuthStatus(userID uint64) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	ConvertUserImageURLs(user)
	return user, nil
}

func GetUserByID(userID uint64) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	ConvertUserImageURLs(user)
	return user, nil
}

func CheckUserAuthStatus(userID uint64) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.AuthStatus != 2 {
		return util.NewBizError(util.ErrCodeUserNotAuth, "用户未通过认证")
	}

	if user.Status != 1 {
		return util.NewBizError(util.ErrCodeUserDisabled, "用户已被禁用")
	}

	return nil
}

func GetUserByPhone(phone string) (*model.User, error) {
	user, err := repository.GetUserByPhone(phone)
	if err != nil {
		return nil, err
	}
	ConvertUserImageURLs(user)
	return user, nil
}
