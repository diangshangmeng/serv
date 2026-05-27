package repository

import (
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreateCoupon(coupon *model.Coupon) error {
	return config.DB.Create(coupon).Error
}

func GetCouponByCode(code string) (*model.Coupon, error) {
	var coupon model.Coupon
	err := config.DB.Where("code = ?", code).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func GetCouponByID(id uint64) (*model.Coupon, error) {
	var coupon model.Coupon
	err := config.DB.Where("id = ?", id).First(&coupon).Error
	if err != nil {
		return nil, err
	}
	return &coupon, nil
}

func GetCouponsByUserID(userID uint64) ([]model.Coupon, error) {
	var coupons []model.Coupon
	err := config.DB.Where("user_id = ?", userID).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

func GetAvailableCouponsByUserID(userID uint64) ([]model.Coupon, error) {
	var coupons []model.Coupon
	err := config.DB.Where("user_id = ? AND status = ?", userID, 0).Find(&coupons).Error
	if err != nil {
		return nil, err
	}
	return coupons, nil
}

func UpdateCoupon(coupon *model.Coupon) error {
	return config.DB.Save(coupon).Error
}

func UpdateCouponStatus(couponID uint64, status int) error {
	return config.DB.Model(&model.Coupon{}).Where("id = ?", couponID).Update("status", status).Error
}

type CouponListResult struct {
	Coupons  []model.Coupon `json:"coupons"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

func GetCouponList(page, pageSize int, code string, status *int) (*CouponListResult, error) {
	var coupons []model.Coupon
	var total int64

	query := config.DB.Model(&model.Coupon{})

	if code != "" {
		query = query.Where("code = ?", code)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&coupons).Error; err != nil {
		return nil, err
	}

	return &CouponListResult{
		Coupons:  coupons,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

type CouponDetailResult struct {
	Coupon   model.Coupon    `json:"coupon"`
	User     model.User      `json:"user"`
	Listing  *model.Listing  `json:"listing,omitempty"`
	Order    *model.Order    `json:"order,omitempty"`
}

func GetCouponDetail(couponID uint64) (*CouponDetailResult, error) {
	var coupon model.Coupon
	if err := config.DB.Where("id = ?", couponID).First(&coupon).Error; err != nil {
		return nil, err
	}

	var user model.User
	if err := config.DB.Where("id = ?", coupon.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	result := &CouponDetailResult{
		Coupon: coupon,
		User:   user,
	}

	if coupon.Status == 2 {
		var listing model.Listing
		if err := config.DB.Where("coupon_id = ?", couponID).First(&listing).Error; err == nil {
			result.Listing = &listing

			var order model.Order
			if err := config.DB.Where("listing_id = ?", listing.ID).First(&order).Error; err == nil {
				result.Order = &order
			}
		}
	}

	return result, nil
}