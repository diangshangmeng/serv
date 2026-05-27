package service

import (
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func CreateListing(userID uint64, couponID uint64, sellingPrice float64) error {
	err := CheckUserAuthStatus(userID)
	if err != nil {
		return err
	}

	coupon, err := repository.GetCouponByID(couponID)
	if err != nil {
		return err
	}

	if coupon.UserID != uint(userID) {
		return util.NewBizError(util.ErrCodeNoPermission, "无权操作此券码")
	}

	if coupon.Status != 0 {
		return util.NewBizError(util.ErrCodeCouponStatusError, "券码状态不正确")
	}

	listing := &model.Listing{
		CouponID:     uint(couponID),
		SellingPrice: sellingPrice,
		Status:       0,
	}

	err = repository.CreateListing(listing)
	if err != nil {
		return err
	}

	coupon.Status = 1
	return repository.UpdateCoupon(coupon)
}

func GetListingList(userID uint64) ([]model.Listing, error) {
	err := CheckUserAuthStatus(userID)
	if err != nil {
		return nil, err
	}

	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	return repository.GetListingsByCity(user.CityID)
}

func GetMyListings(userID uint64) ([]model.Listing, error) {
	err := CheckUserAuthStatus(userID)
	if err != nil {
		return nil, err
	}

	return repository.GetListingsByUserID(userID)
}