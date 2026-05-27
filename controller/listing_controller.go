package controller

import (
	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

func CreateListing(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	var req struct {
		CouponID     uint64  `json:"coupon_id" binding:"required"`
		SellingPrice float64 `json:"selling_price" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.CreateListing(userID, req.CouponID, req.SellingPrice)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "挂单成功")
}

func GetListingList(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	listings, err := service.GetListingList(userID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, listings)
}

func GetMyListings(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	listings, err := service.GetMyListings(userID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, listings)
}
