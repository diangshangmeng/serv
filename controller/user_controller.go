package controller

import (
	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

func SubmitAuth(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	var req struct {
		CityID          uint   `json:"city_id" binding:"required"`
		IDCardFront     string `json:"id_card_front" binding:"required"`
		IDCardBack      string `json:"id_card_back" binding:"required"`
		BusinessLicense string `json:"business_license" binding:"required"`
		PaymentCode     string `json:"payment_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.SubmitAuth(userID, req.CityID, req.IDCardFront, req.IDCardBack, req.BusinessLicense, req.PaymentCode)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "提交成功，等待审核")
}

func GetAuthStatus(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	user, err := service.GetAuthStatus(userID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"authStatus":   user.AuthStatus,
		"auditRemark":  user.AuditRemark,
		"status":       user.Status,
	}

	util.ResponseSuccess(c, result)
}
