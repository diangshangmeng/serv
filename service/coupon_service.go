package service

import (
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func AddCoupon(userID uint64, code string, originalPrice float64) error {
	err := CheckUserAuthStatus(userID)
	if err != nil {
		return err
	}

	_, err = repository.GetCouponByCode(code)
	if err == nil {
		return util.NewBizError(util.ErrCodeCouponExists, "券码已存在")
	}

	coupon := &model.Coupon{
		UserID:        uint(userID),
		Code:          code,
		OriginalPrice: originalPrice,
		Status:        0,
	}

	return repository.CreateCoupon(coupon)
}

func GetMyCoupons(userID uint64) ([]model.Coupon, error) {
	err := CheckUserAuthStatus(userID)
	if err != nil {
		return nil, err
	}

	return repository.GetAvailableCouponsByUserID(userID)
}

type CouponListResponse struct {
	List     []CouponInfo `json:"list"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

type CouponInfo struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	Code          string  `json:"code"`
	OriginalPrice float64 `json:"original_price"`
	Status        int     `json:"status"`
	CreatedAt     string  `json:"created_at"`
	UserPhone     string  `json:"user_phone"`
	UserCityID    uint    `json:"user_city_id"`
}

func GetCouponList(page, pageSize int, code string, status *int) (*CouponListResponse, error) {
	result, err := repository.GetCouponList(page, pageSize, code, status)
	if err != nil {
		return nil, err
	}

	list := make([]CouponInfo, 0, len(result.Coupons))
	for _, c := range result.Coupons {
		user, _ := repository.GetUserByID(uint64(c.UserID))
		info := CouponInfo{
			ID:            c.ID,
			UserID:        c.UserID,
			Code:          c.Code,
			OriginalPrice: c.OriginalPrice,
			Status:        c.Status,
			CreatedAt:     c.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if user != nil {
			info.UserPhone = user.Phone
			info.UserCityID = user.CityID
		}
		list = append(list, info)
	}

	return &CouponListResponse{
		List:     list,
		Total:    result.Total,
		Page:     result.Page,
		PageSize: result.PageSize,
	}, nil
}

type CouponDetailResponse struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	Code          string  `json:"code"`
	OriginalPrice float64 `json:"original_price"`
	Status        int     `json:"status"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	User          struct {
		ID       uint   `json:"id"`
		Phone    string `json:"phone"`
		CityID   uint   `json:"city_id"`
		AuthStatus int  `json:"auth_status"`
	} `json:"user"`
	Listing *ListingInfo `json:"listing,omitempty"`
	Order   *OrderInfo   `json:"order,omitempty"`
}

type ListingInfo struct {
	ID          uint    `json:"id"`
	SellingPrice float64 `json:"selling_price"`
	Status      int     `json:"status"`
	CreatedAt   string  `json:"created_at"`
}

type OrderInfo struct {
	ID        uint   `json:"id"`
	OrderNo   string `json:"order_no"`
	BuyerID   uint   `json:"buyer_id"`
	SellerID  uint   `json:"seller_id"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

func GetCouponDetail(couponID uint64) (*CouponDetailResponse, error) {
	detail, err := repository.GetCouponDetail(couponID)
	if err != nil {
		return nil, util.NewBizError(util.ErrCodeCouponNotFound, "券码不存在")
	}

	resp := &CouponDetailResponse{
		ID:            detail.Coupon.ID,
		UserID:        detail.Coupon.UserID,
		Code:          detail.Coupon.Code,
		OriginalPrice: detail.Coupon.OriginalPrice,
		Status:        detail.Coupon.Status,
		CreatedAt:     detail.Coupon.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     detail.Coupon.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	resp.User.ID = detail.User.ID
	resp.User.Phone = detail.User.Phone
	resp.User.CityID = detail.User.CityID
	resp.User.AuthStatus = detail.User.AuthStatus

	if detail.Listing != nil {
		resp.Listing = &ListingInfo{
			ID:           detail.Listing.ID,
			SellingPrice: detail.Listing.SellingPrice,
			Status:       detail.Listing.Status,
			CreatedAt:    detail.Listing.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	if detail.Order != nil {
		resp.Order = &OrderInfo{
			ID:        detail.Order.ID,
			OrderNo:   detail.Order.OrderNo,
			BuyerID:   detail.Order.BuyerID,
			SellerID:  detail.Order.SellerID,
			Status:    detail.Order.Status,
			CreatedAt: detail.Order.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return resp, nil
}