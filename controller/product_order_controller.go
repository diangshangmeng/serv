package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

func CreateProductOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)
	phone := c.MustGet("phone").(string)

	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	order, _, err := service.CreateProductOrder(req.ProductID, userID, phone)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"order": order,
		"seller_payment_code_image_url": order.PaymentCodeImageURL,
	}

	util.ResponseSuccessWithMessage(c, result, "创建订单成功，请在15分钟内完成支付")
}

func PlaceOrder(c *gin.Context) {
	productIDStr := c.Query("product_id")
	if productIDStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	userID := c.MustGet("userID").(uint64)

	if err := service.PlaceOrder(productID, userID, ""); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message":     "下单成功，请在15分钟内完成支付",
		"product_id": productID,
	})
}

func CompletePayment(c *gin.Context) {
	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if err := service.CompleteProductPayment(req.ProductID); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message":     "支付成功，等待卖家确认",
		"product_id": req.ProductID,
	})
}

func ConfirmProduct(c *gin.Context) {
	productIDStr := c.Query("product_id")
	if productIDStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if err := service.ConfirmProductOrder(productID); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message":     "确认成功，交易完成",
		"product_id": productID,
	})
}

func CancelProduct(c *gin.Context) {
	var req struct {
		ProductID uint64 `json:"product_id" binding:"required"`
		Reason    string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if req.Reason == "" {
		req.Reason = "用户取消"
	}

	if err := service.CancelProductOrder(req.ProductID, req.Reason); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message":     "订单已取消",
		"product_id": req.ProductID,
	})
}

func GetOrderStatus(c *gin.Context) {
	productIDStr := c.Query("product_id")
	if productIDStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	product, err := service.GetProductByID(productID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	response := gin.H{
		"product_id": productID,
		"status":     product.Status,
		"status_text": getStatusText(product.Status),
	}

	if product.Status == 2 {
		order, err := service.GetActiveOrderByProductID(productID)
		if err == nil && order != nil && order.OrderTime != nil {
			orderTime := order.OrderTime.Add(15 * time.Minute)
			response["payment_deadline"] = orderTime
			response["payment_remaining_seconds"] = int(time.Until(orderTime).Seconds())
			if time.Until(orderTime) < 0 {
				response["payment_expired"] = true
			} else {
				response["payment_expired"] = false
			}
		}
	}

	if product.PayTime != nil {
		confirmTime := product.PayTime.Add(30 * time.Minute)
		response["confirm_deadline"] = confirmTime
		response["confirm_remaining_seconds"] = int(time.Until(confirmTime).Seconds())
		if time.Until(confirmTime) < 0 {
			response["confirm_expired"] = true
		} else {
			response["confirm_expired"] = false
		}
	}

	if product.LockReason != "" {
		response["lock_reason"] = product.LockReason
	}

	util.ResponseSuccess(c, response)
}

func getStatusText(status int) string {
	switch status {
	case 0:
		return "可购买"
	case 1:
		return "等待支付"
	case 2:
		return "已锁定"
	case 4:
		return "等待确认"
	default:
		return "未知状态"
	}
}
