package controller

import (
	"strconv"
	"time"

	"go.uber.org/zap"
	"voucher-platform/service"
	"voucher-platform/util"

	"github.com/gin-gonic/gin"
)

type OrderListItem struct {
	OrderNo           string `json:"order_no"`
	ProductTitle      string `json:"product_title"`
	ProductImageURL   string `json:"product_image_url"`
	Price             int64  `json:"price"`
	Status            int    `json:"status"`
	StatusText        string `json:"status_text"`
	OrderType         string `json:"order_type"`
	CreatedAt         string `json:"created_at"`
	ExpiredAt         string `json:"expired_at"`
	CounterpartyPhone string `json:"counterparty_phone"`
	PaymentVoucher    string `json:"payment_voucher"`
}

type OrderDetailItem struct {
	ID                  uint   `json:"id"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	OrderNo             string `json:"order_no"`
	ProductID           uint   `json:"product_id"`
	ProductTitle        string `json:"product_title"`
	ProductDescription  string `json:"product_description"`
	SellerID            uint   `json:"seller_id"`
	SellerPhone         string `json:"seller_phone"`
	ProductImageURL     string `json:"product_image_url"`
	PaymentCodeImageURL string `json:"payment_code_image_url"`
	BuyerID             uint   `json:"buyer_id"`
	BuyerPhone          string `json:"buyer_phone"`
	PaymentVoucher      string `json:"payment_voucher"`
	Status              int    `json:"status"`
	StatusText          string `json:"status_text"`
	OrderType           string `json:"order_type"`
	Price               int64  `json:"price"`
	ExpiredAt           string `json:"expired_at"`
	CounterpartyPhone   string `json:"counterparty_phone"`
}

func GetStatusText(status int) string {
	switch status {
	case 0:
		return "未付款"
	case 1:
		return "已付款"
	case 2:
		return "已完成"
	case 3:
		return "已取消"
	case 4:
		return "超时关闭"
	default:
		return "未知状态"
	}
}

func MaskPhone(phone string) string {
	if len(phone) < 7 {
		return phone
	}
	return phone[:3] + "****" + phone[len(phone)-4:]
}

func CreateOrder(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)

	var req struct {
		ListingID uint64 `json:"listing_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	order, seller, err := service.CreateOrder(userID, req.ListingID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"order": gin.H{
			"order_no":               order.OrderNo,
			"product_title":          order.ProductTitle,
			"product_description":    order.ProductDescription,
			"product_image_url":      order.ProductImageURL,
			"payment_code_image_url": order.PaymentCodeImageURL,
			"price":                  order.Price,
			"expired_at":             order.ExpiredAt.Format(time.RFC3339),
		},
		"seller": gin.H{
			"phone": seller.Phone,
		},
	}

	util.ResponseSuccess(c, result)
}

func UploadVoucher(c *gin.Context) {
	util.GetLogger().Info("UploadVoucher - 开始处理请求")

	userID := c.MustGet("userID").(uint64)

	var req struct {
		OrderNo        string `json:"order_no" binding:"required"`
		PaymentVoucher string `json:"payment_voucher"`
		VoucherURL     string `json:"voucher_url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.GetLogger().Error("UploadVoucher - 绑定参数失败", zap.Error(err))
		util.ResponseError(c, 400, "参数错误")
		return
	}

	var voucherURL string
	var usedField string
	if req.PaymentVoucher != "" {
		voucherURL = req.PaymentVoucher
		usedField = "payment_voucher"
	} else if req.VoucherURL != "" {
		voucherURL = req.VoucherURL
		usedField = "voucher_url"
	} else {
		util.GetLogger().Error("UploadVoucher - 缺少凭证URL参数",
			util.StringField("order_no", req.OrderNo),
			util.Uint64Field("user_id", userID))
		util.ResponseError(c, 400, "缺少凭证URL参数")
		return
	}

	util.GetLogger().Info("UploadVoucher - 参数绑定成功",
		util.StringField("order_no", req.OrderNo),
		util.Uint64Field("user_id", userID),
		util.StringField("used_field", usedField),
		util.StringField("voucher_url", voucherURL))

	err := service.UploadVoucher(req.OrderNo, voucherURL, userID)
	if err != nil {
		util.GetLogger().Error("UploadVoucher - 上传凭证失败", zap.Error(err),
			util.StringField("order_no", req.OrderNo),
			util.Uint64Field("user_id", userID))
		util.ResponseBizError(c, err)
		return
	}

	util.GetLogger().Info("UploadVoucher - 上传凭证成功",
		util.StringField("order_no", req.OrderNo),
		util.Uint64Field("user_id", userID))
	util.ResponseSuccess(c, nil)
}

func ConfirmOrder(c *gin.Context) {
	phone := c.MustGet("phone").(string)

	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.ConfirmOrder(req.OrderNo, phone)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, nil)
}

func SellerCancelOrder(c *gin.Context) {
	util.GetLogger().Info("SellerCancelOrder - 开始处理请求")

	phone := c.MustGet("phone").(string)

	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.GetLogger().Error("SellerCancelOrder - 绑定参数失败", zap.Error(err))
		util.ResponseError(c, 400, "参数错误")
		return
	}

	util.GetLogger().Info("SellerCancelOrder - 参数绑定成功",
		util.StringField("order_no", req.OrderNo),
		util.StringField("phone", phone))

	err := service.SellerCancelOrder(req.OrderNo, phone)
	if err != nil {
		util.GetLogger().Error("SellerCancelOrder - 取消订单失败", zap.Error(err),
			util.StringField("order_no", req.OrderNo),
			util.StringField("phone", phone))
		util.ResponseBizError(c, err)
		return
	}

	util.GetLogger().Info("SellerCancelOrder - 取消订单成功",
		util.StringField("order_no", req.OrderNo),
		util.StringField("phone", phone))
	util.ResponseSuccess(c, nil)
}

func GetMyBuyOrders(c *gin.Context) {
	phone := c.MustGet("phone").(string)

	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	orders, total, err := service.GetMyBuyOrders(phone, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderItems := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		orderItems = append(orderItems, OrderListItem{
			OrderNo:           order.OrderNo,
			ProductTitle:      order.ProductTitle,
			ProductImageURL:   order.ProductImageURL,
			Price:             order.Price,
			Status:            order.Status,
			StatusText:        GetStatusText(order.Status),
			OrderType:         "buy",
			CreatedAt:         order.CreatedAt.Format(time.RFC3339),
			ExpiredAt:         order.ExpiredAt.Format(time.RFC3339),
			CounterpartyPhone: MaskPhone(order.SellerPhone),
			PaymentVoucher:    order.PaymentVoucher,
		})
	}

	result := gin.H{
		"orders":    orderItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetMySellOrders(c *gin.Context) {
	phone := c.MustGet("phone").(string)

	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	orders, total, err := service.GetMySellOrders(phone, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderItems := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		orderItems = append(orderItems, OrderListItem{
			OrderNo:           order.OrderNo,
			ProductTitle:      order.ProductTitle,
			ProductImageURL:   order.ProductImageURL,
			Price:             order.Price,
			Status:            order.Status,
			StatusText:        GetStatusText(order.Status),
			OrderType:         "sell",
			CreatedAt:         order.CreatedAt.Format(time.RFC3339),
			ExpiredAt:         order.ExpiredAt.Format(time.RFC3339),
			CounterpartyPhone: MaskPhone(order.BuyerPhone),
			PaymentVoucher:    order.PaymentVoucher,
		})
	}

	result := gin.H{
		"orders":    orderItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetOrderDetail(c *gin.Context) {
	userID := c.MustGet("userID").(uint64)
	orderNo := c.Query("order_no")

	if orderNo == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	order, err := service.GetOrderDetail(orderNo, userID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	var orderType string
	var counterpartyPhone string

	if uint64(order.BuyerID) == userID {
		orderType = "buy"
		counterpartyPhone = MaskPhone(order.SellerPhone)
	} else {
		orderType = "sell"
		counterpartyPhone = MaskPhone(order.BuyerPhone)
	}

	orderDetail := OrderDetailItem{
		ID:                  order.ID,
		CreatedAt:           order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           order.UpdatedAt.Format(time.RFC3339),
		OrderNo:             order.OrderNo,
		ProductID:           order.ProductID,
		ProductTitle:        order.ProductTitle,
		ProductDescription:  order.ProductDescription,
		SellerID:            order.SellerID,
		SellerPhone:         order.SellerPhone,
		ProductImageURL:     order.ProductImageURL,
		PaymentCodeImageURL: order.PaymentCodeImageURL,
		BuyerID:             order.BuyerID,
		BuyerPhone:          order.BuyerPhone,
		PaymentVoucher:      order.PaymentVoucher,
		Status:              order.Status,
		StatusText:          GetStatusText(order.Status),
		OrderType:           orderType,
		Price:               order.Price,
		ExpiredAt:           order.ExpiredAt.Format(time.RFC3339),
		CounterpartyPhone:   counterpartyPhone,
	}

	util.ResponseSuccess(c, orderDetail)
}

func GetAppBuyerOrders(c *gin.Context) {
	phone := c.MustGet("phone").(string)

	page := 1
	pageSize := 20
	var status *int

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	if s := c.Query("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
			status = &parsed
		}
	}

	orders, total, err := service.GetBuyerOrdersByPhone(phone, status, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderItems := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		var orderType string
		var counterpartyPhone string

		if order.BuyerPhone == phone {
			orderType = "buy"
			counterpartyPhone = MaskPhone(order.SellerPhone)
		} else {
			orderType = "sell"
			counterpartyPhone = MaskPhone(order.BuyerPhone)
		}

		orderItems = append(orderItems, OrderListItem{
			OrderNo:           order.OrderNo,
			ProductTitle:      order.ProductTitle,
			ProductImageURL:   order.ProductImageURL,
			Price:             order.Price,
			Status:            order.Status,
			StatusText:        GetStatusText(order.Status),
			OrderType:         orderType,
			CreatedAt:         order.CreatedAt.Format(time.RFC3339),
			ExpiredAt:         order.ExpiredAt.Format(time.RFC3339),
			CounterpartyPhone: counterpartyPhone,
			PaymentVoucher:    order.PaymentVoucher,
		})
	}

	result := gin.H{
		"orders":    orderItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetAppSellerOrders(c *gin.Context) {
	phone := c.MustGet("phone").(string)

	page := 1
	pageSize := 20
	var status *int

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	if s := c.Query("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
			status = &parsed
		}
	}

	orders, total, err := service.GetSellerOrdersByPhone(phone, status, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderItems := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		orderItems = append(orderItems, OrderListItem{
			OrderNo:           order.OrderNo,
			ProductTitle:      order.ProductTitle,
			ProductImageURL:   order.ProductImageURL,
			Price:             order.Price,
			Status:            order.Status,
			StatusText:        GetStatusText(order.Status),
			OrderType:         "sell",
			CreatedAt:         order.CreatedAt.Format(time.RFC3339),
			ExpiredAt:         order.ExpiredAt.Format(time.RFC3339),
			CounterpartyPhone: MaskPhone(order.BuyerPhone),
			PaymentVoucher:    order.PaymentVoucher,
		})
	}

	result := gin.H{
		"orders":    orderItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetAppAllOrders(c *gin.Context) {
	phone := c.MustGet("phone").(string)

	page := 1
	pageSize := 20
	var status *int

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	if s := c.Query("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
			status = &parsed
		}
	}

	orders, total, err := service.GetAllOrdersByPhone(phone, status, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderItems := make([]OrderListItem, 0, len(orders))
	for _, order := range orders {
		var orderType string
		var counterpartyPhone string

		if order.BuyerPhone == phone {
			orderType = "buy"
			counterpartyPhone = MaskPhone(order.SellerPhone)
		} else {
			orderType = "sell"
			counterpartyPhone = MaskPhone(order.BuyerPhone)
		}

		orderItems = append(orderItems, OrderListItem{
			OrderNo:           order.OrderNo,
			ProductTitle:      order.ProductTitle,
			ProductImageURL:   order.ProductImageURL,
			Price:             order.Price,
			Status:            order.Status,
			StatusText:        GetStatusText(order.Status),
			OrderType:         orderType,
			CreatedAt:         order.CreatedAt.Format(time.RFC3339),
			ExpiredAt:         order.ExpiredAt.Format(time.RFC3339),
			CounterpartyPhone: counterpartyPhone,
			PaymentVoucher:    order.PaymentVoucher,
		})
	}

	result := gin.H{
		"orders":    orderItems,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}
