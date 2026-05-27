package controller

import (
	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

func AddCoupon(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	var req struct {
		Code          string  `json:"code" binding:"required"`
		OriginalPrice float64 `json:"original_price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.AddCoupon(userID, req.Code, req.OriginalPrice)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "添加成功")
}

func GetMyCoupons(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	coupons, err := service.GetMyCoupons(userID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, coupons)
}