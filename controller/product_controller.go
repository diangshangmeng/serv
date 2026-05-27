package controller

import (
	"log"
	"strconv"

	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/service"
	"voucher-platform/util"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	adminID := c.MustGet("adminID").(uint64)
	adminPhone := c.MustGet("adminPhone").(string)

	var req struct {
		Title           string `json:"title" binding:"required"`
		Description     string `json:"description"`
		Price           int64  `json:"price" binding:"required"`
		DisplayImageID  uint   `json:"display_image_id"`
		PaymentImageID  uint   `json:"payment_image_id"`
		DisplayImageURL string `json:"display_image_url"`
		PaymentImageURL string `json:"payment_image_url"`
		CityID          uint   `json:"city_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if req.CityID != 0 {
		if _, err := repository.GetCityByID(uint64(req.CityID)); err != nil {
			util.ResponseError(c, 400, "城市不存在")
			return
		}
	}

	product := &model.Product{
		Title:           req.Title,
		Description:     req.Description,
		Price:           req.Price,
		DisplayImageID:  req.DisplayImageID,
		PaymentImageID:  req.PaymentImageID,
		DisplayImageURL: req.DisplayImageURL,
		PaymentImageURL: req.PaymentImageURL,
		CityID:          req.CityID,
	}

	if err := service.CreateProduct(product, uint(adminID), adminPhone); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func GetProductList(c *gin.Context) {
	page := 1
	pageSize := 1000
	var cityID *uint = nil

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

	if cid := c.Query("city_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			cityUint := uint(parsed)
			cityID = &cityUint
		}
	}

	list, total, err := service.GetProductList(page, pageSize, cityID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}
func GetProductListForDashboard(c *gin.Context) {
	page := 1
	pageSize := 1000
	var cityID *uint = nil

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

	if cid := c.Query("city_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			cityUint := uint(parsed)
			cityID = &cityUint
		}
	}

	list, total, err := service.GetProductListForDashboard(page, pageSize, cityID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	var req struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		Price           int64  `json:"price"`
		DisplayImageID  uint   `json:"display_image_id"`
		PaymentImageID  uint   `json:"payment_image_id"`
		DisplayImageURL string `json:"display_image_url"`
		PaymentImageURL string `json:"payment_image_url"`
		CityID          uint   `json:"city_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if req.CityID != 0 {
		if _, err := repository.GetCityByID(uint64(req.CityID)); err != nil {
			util.ResponseError(c, 400, "城市不存在")
			return
		}
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	if req.Title != "" {
		product.Title = req.Title
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price != 0 {
		product.Price = req.Price
	}
	if req.DisplayImageID != 0 {
		product.DisplayImageID = req.DisplayImageID
	}
	if req.PaymentImageID != 0 {
		product.PaymentImageID = req.PaymentImageID
	}
	if req.DisplayImageURL != "" {
		product.DisplayImageURL = req.DisplayImageURL
	}
	if req.PaymentImageURL != "" {
		product.PaymentImageURL = req.PaymentImageURL
	}
	if req.CityID != 0 {
		product.CityID = req.CityID
	}

	if err := service.UpdateProduct(product); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	adminPhone := c.MustGet("adminPhone").(string)

	if err := service.DeleteProduct(id, adminPhone); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, nil)
}

func GetAppProductList(c *gin.Context) {
	page := 1
	pageSize := 1000
	var cityID *uint = nil

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

	if cid := c.Query("city_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			cityUint := uint(parsed)
			cityID = &cityUint
		}
	}

	list, total, err := service.GetProductList(page, pageSize, cityID)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetAppProductDetail(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func PublishProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if err := service.PublishProduct(id); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func UnpublishProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if req.Reason == "" {
		req.Reason = "管理员下架"
	}

	if err := service.UnpublishProduct(id, req.Reason); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func UnlockProduct(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if err := service.UnlockProduct(id); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, product)
}

func GetMyProducts(c *gin.Context) {
	log.Println("[GetMyProducts] 开始处理请求")

	phone := c.MustGet("phone").(string)
	log.Printf("[GetMyProducts] 从JWT获取phone成功，phone=%s", phone)

	page := 1
	pageSize := 1000
	var cityID *uint = nil

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

	if cid := c.Query("city_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			cityUint := uint(parsed)
			cityID = &cityUint
		}
	}

	log.Printf("[GetMyProducts] 查询参数，page=%d, pageSize=%d", page, pageSize)

	list, total, err := service.GetMyProducts(phone, page, pageSize, cityID)
	if err != nil {
		log.Printf("[GetMyProducts] 查询失败，error=%v", err)
		util.ResponseBizError(c, err)
		return
	}

	log.Printf("[GetMyProducts] 查询成功，total=%d", total)

	result := gin.H{
		"list":      list,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func AppPublishProduct(c *gin.Context) {
	log.Println("[AppPublishProduct] 开始处理请求")

	idStr := c.Param("id")
	if idStr == "" {
		log.Println("[AppPublishProduct] 参数错误，id为空")
		util.ResponseError(c, 400, "参数错误")
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Printf("[AppPublishProduct] 参数错误，id格式错误，idStr=%s", idStr)
		util.ResponseError(c, 400, "参数错误")
		return
	}

	phone := c.MustGet("phone").(string)
	log.Printf("[AppPublishProduct] 从JWT获取phone成功，phone=%s, productID=%d", phone, id)

	if err := service.AppPublishProduct(id, phone); err != nil {
		log.Printf("[AppPublishProduct] 失败，error=%v", err)
		util.ResponseBizError(c, err)
		return
	}

	product, err := service.GetProductByID(id)
	if err != nil {
		log.Printf("[AppPublishProduct] 获取商品失败，error=%v", err)
		util.ResponseBizError(c, err)
		return
	}

	log.Printf("[AppPublishProduct] 成功，productID=%d", id)
	util.ResponseSuccess(c, product)
}
