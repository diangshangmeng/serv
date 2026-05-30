package controller

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"voucher-platform/service"
	"voucher-platform/util"
)

type AdminOrderListItem struct {
	ID                  uint   `json:"id"`
	OrderNo             string `json:"order_no"`
	ProductTitle        string `json:"product_title"`
	SellerPhone         string `json:"seller_phone"`
	BuyerPhone          string `json:"buyer_phone"`
	Status              int    `json:"status"`
	StatusText          string `json:"status_text"`
	Price               int64  `json:"price"`
	CreatedAt           string `json:"created_at"`
	ExpiredAt           string `json:"expired_at"`
	PaymentVoucher      string `json:"payment_voucher"`
}

type AdminOrderDetailItem struct {
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
	Price               int64  `json:"price"`
	ExpiredAt           string `json:"expired_at"`
}

func AdminLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	token, err := service.AdminLogin(req.Username, req.Password)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"token": token,
	}

	util.ResponseSuccess(c, result)
}

func UpdateAdminPassword(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint64)

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if err := service.UpdateAdminPassword(adminID, req.OldPassword, req.NewPassword); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "密码修改成功")
}

func GetPendingUsers(c *gin.Context) {
	users, err := service.GetPendingUsers()
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, users)
}

func AuditUser(c *gin.Context) {
	type AuditRequest struct {
		UserID interface{} `json:"user_id"`
		Pass   interface{} `json:"pass"`
		Remark string      `json:"remark"`
	}

	var req AuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "参数绑定错误",
			"error":   err.Error(),
		})
		return
	}

	var userID uint64
	switch v := req.UserID.(type) {
	case float64:
		userID = uint64(v)
	case string:
		if u, err := strconv.ParseUint(v, 10, 64); err == nil {
			userID = u
		}
	case int:
		userID = uint64(v)
	case int64:
		userID = uint64(v)
	}

	var pass bool
	switch v := req.Pass.(type) {
	case bool:
		pass = v
	case string:
		pass = v == "true"
	}

	if userID == 0 {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "user_id 无效或未提供",
			"req":     req,
		})
		return
	}

	err := service.AuditUser(userID, pass, req.Remark)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "审核完成")
}

func GetAllUsers(c *gin.Context) {
	users, err := service.GetAllUsers()
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, users)
}

func SetUserStatus(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
		Status int    `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.SetUserStatus(req.UserID, req.Status)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "操作成功")
}

func ResetUserPassword(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	newPassword, err := service.ResetUserPassword(req.UserID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"password": newPassword,
	}

	util.ResponseSuccess(c, result)
}

func GetAllOrders(c *gin.Context) {
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

	orderNo := c.Query("order_no")
	phone := c.Query("phone")
	
	var status *int
	if s := c.Query("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed >= 0 {
			status = &parsed
		}
	}

	var startTime, endTime *time.Time
	if st := c.Query("start_time"); st != "" {
		if t, err := time.Parse(time.RFC3339, st); err == nil {
			startTime = &t
		}
	}
	if et := c.Query("end_time"); et != "" {
		if t, err := time.Parse(time.RFC3339, et); err == nil {
			endTime = &t
		}
	}

	orders, total, err := service.AdminGetOrders(orderNo, phone, status, startTime, endTime, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderItems := make([]AdminOrderListItem, 0, len(orders))
	for _, order := range orders {
		orderItems = append(orderItems, AdminOrderListItem{
			ID:             order.ID,
			OrderNo:        order.OrderNo,
			ProductTitle:   order.ProductTitle,
			SellerPhone:    order.SellerPhone,
			BuyerPhone:     order.BuyerPhone,
			Status:         order.Status,
			StatusText:     GetStatusText(order.Status),
			Price:          order.Price,
			CreatedAt:      order.CreatedAt.Format(time.RFC3339),
			ExpiredAt:      order.ExpiredAt.Format(time.RFC3339),
			PaymentVoucher: order.PaymentVoucher,
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

func AdminGetOrderDetail(c *gin.Context) {
	orderNo := c.Query("order_no")
	if orderNo == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	order, err := service.GetOrderByOrderNo(orderNo)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	orderDetail := AdminOrderDetailItem{
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
		Price:               order.Price,
		ExpiredAt:           order.ExpiredAt.Format(time.RFC3339),
	}

	util.ResponseSuccess(c, orderDetail)
}

func GetOrderStats(c *gin.Context) {
	stats, err := service.GetOrderStats()
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, stats)
}

func AdminCancelOrder(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.CancelOrder(req.OrderNo)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "取消成功")
}

func GetCouponList(c *gin.Context) {
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

	code := c.Query("code")

	var status *int
	if s := c.Query("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			status = &parsed
		}
	}

	result, err := service.GetCouponList(page, pageSize, code, status)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, result)
}

func GetCouponDetail(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	result, err := service.GetCouponDetail(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, result)
}

func GetUserDetail(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	user, err := service.GetUserDetail(userID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, user)
}

func SearchUsers(c *gin.Context) {
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

	phone := c.Query("phone")

	var status *int
	if s := c.Query("status"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil {
			status = &parsed
		}
	}

	users, total, err := service.SearchUsers(phone, status, page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func BanUser(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.BanUser(req.UserID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "封禁成功")
}

func UnbanUser(c *gin.Context) {
	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.UnbanUser(req.UserID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "解封成功")
}

func UpdateUser(c *gin.Context) {
	var req struct {
		UserID          uint64  `json:"user_id" binding:"required"`
		CityID          *uint   `json:"city_id"`
		Status          *int    `json:"status"`
		AuthStatus      *int    `json:"auth_status"`
		AuditRemark     string  `json:"audit_remark"`
		IDCardFront     string  `json:"id_card_front"`
		IDCardBack      string  `json:"id_card_back"`
		BusinessLicense string  `json:"business_license"`
		PaymentCode     string  `json:"payment_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	user, err := service.UpdateUser(req.UserID, req.CityID, req.Status, req.AuthStatus, req.AuditRemark, req.IDCardFront, req.IDCardBack, req.BusinessLicense, req.PaymentCode)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, user)
}

func AdminSellerCancelOrder(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.AdminSellerCancelOrder(req.OrderNo)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "取消成功")
}

func AdminSellerConfirmOrder(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	err := service.AdminSellerConfirmOrder(req.OrderNo)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseMessage(c, "确认成功")
}
