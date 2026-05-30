package util

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrCodeSuccess           = 0
	ErrCodeSystem            = 500
	ErrCodeParamInvalid      = 400
	ErrCodeUnauthorized      = 401
	ErrCodeForbidden         = 403
	ErrCodeNotFound          = 404

	ErrCodeUserExists        = 1001
	ErrCodeUserNotFound      = 1002
	ErrCodePasswordError     = 1003
	ErrCodeUserDisabled      = 1004
	ErrCodeUserNotAuth       = 1005

	ErrCodeCouponExists      = 2001
	ErrCodeCouponNotFound    = 2002
	ErrCodeCouponStatusError = 2003

	ErrCodeListingNotFound  = 3001
	ErrCodeListingStatusError= 3002
	ErrCodeCannotBuyOwn      = 3003
	ErrCodeNotSameCity       = 3004

	ErrCodeOrderNotFound     = 4001
	ErrCodeOrderStatusError  = 4002
	ErrCodeNoPermission      = 4003

	ErrCodeAdminAuthFailed   = 5001
	ErrCodeAdminNotFound     = 5002

	ErrCodeImageNotFound     = 6001
	ErrCodeImageAlreadyUsed  = 6002
	ErrCodeProductNotFound     = 7001
	ErrCodeProductNoPermission = 7002
	ErrCodeProductUnavailable  = 7003
	ErrCodeProductStatusError = 7004
	ErrCodePaymentCodeImageNotFound = 8001
)

type BizError struct {
	Code    int
	Message string
	Err     error
}

func (e *BizError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *BizError) Unwrap() error {
	return e.Err
}

func NewBizError(code int, message string) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
	}
}

func WrapError(code int, message string, err error) *BizError {
	return &BizError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func IsBizError(err error) bool {
	_, ok := err.(*BizError)
	return ok
}

func GetBizError(err error) (*BizError, bool) {
	bizErr, ok := err.(*BizError)
	return bizErr, ok
}

func ResponseBizError(c *gin.Context, err error) {
	if bizErr, ok := GetBizError(err); ok {
		statusCode := getHttpStatus(bizErr.Code)
		c.JSON(statusCode, Response{
			Code:    bizErr.Code,
			Message: bizErr.Message,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		Code:    ErrCodeSystem,
		Message: "系统错误",
		Data:    nil,
	})
}

func getHttpStatus(errCode int) int {
	switch {
	case errCode >= 5000 && errCode < 6000:
		return http.StatusUnauthorized
	case errCode >= 4000 && errCode < 5000:
		return http.StatusForbidden
	case errCode >= 3000 && errCode < 4000:
		return http.StatusBadRequest
	case errCode >= 2000 && errCode < 3000:
		return http.StatusBadRequest
	case errCode >= 1000 && errCode < 2000:
		return http.StatusBadRequest
	case errCode == ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case errCode == ErrCodeForbidden:
		return http.StatusForbidden
	case errCode == ErrCodeNotFound:
		return http.StatusNotFound
	case errCode == ErrCodeParamInvalid:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
