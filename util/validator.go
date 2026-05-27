package util

import (
	"fmt"
	"reflect"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type ValidationErrors struct {
	Errors map[string]string `json:"errors"`
}

func (v *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %v", v.Errors)
}

var phoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func RegisterCustomValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("phone", validatePhone); err != nil {
			return err
		}
		if err := v.RegisterValidation("positive_amount", validatePositiveAmount); err != nil {
			return err
		}
		if err := v.RegisterValidation("price_range", validatePriceRange); err != nil {
			return err
		}
		if err := v.RegisterValidation("password_strength", validatePasswordStrength); err != nil {
			return err
		}
	}
	return nil
}

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return phoneRegex.MatchString(phone)
}

func validatePositiveAmount(fl validator.FieldLevel) bool {
	if fl.Field().Kind() == reflect.Float64 {
		amount := fl.Field().Float()
		return amount > 0
	}
	if fl.Field().Kind() == reflect.Float32 {
		amount := fl.Field().Float()
		return amount > 0
	}
	return false
}

func validatePriceRange(fl validator.FieldLevel) bool {
	amount := fl.Field().Float()
	return amount >= 0.01 && amount <= 1000000
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 || len(password) > 32 {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required,min=6,max=32,password_strength"`
}

type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,phone"`
	Password string `json:"password" binding:"required"`
}

type AddCouponRequest struct {
	Code          string  `json:"code" binding:"required,min=1,max=50"`
	OriginalPrice float64 `json:"original_price" binding:"required,positive_amount,price_range"`
}

type CreateListingRequest struct {
	CouponID     uint64  `json:"coupon_id" binding:"required,gt=0"`
	SellingPrice float64 `json:"selling_price" binding:"required,positive_amount,price_range"`
}

type CreateOrderRequest struct {
	ListingID uint64 `json:"listing_id" binding:"required,gt=0"`
}

type UploadVoucherRequest struct {
	OrderNo        string `json:"order_no" binding:"required,min=1"`
	PaymentVoucher string `json:"payment_voucher" binding:"required,min=1"`
}

type SubmitAuthRequest struct {
	CityID          uint   `json:"city_id" binding:"required,gt=0"`
	IDCardFront     string `json:"id_card_front" binding:"required,url"`
	IDCardBack      string `json:"id_card_back" binding:"required,url"`
	BusinessLicense string `json:"business_license" binding:"required,url"`
	PaymentCode     string `json:"payment_code" binding:"required,url"`
}

func ValidateRequest(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}
	return nil
}
