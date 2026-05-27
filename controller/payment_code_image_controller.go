package controller

import (
	"strconv"

	"voucher-platform/model"
	"voucher-platform/service"
	"voucher-platform/util"

	"github.com/gin-gonic/gin"
)

func UploadPaymentCodeImage(c *gin.Context) {
	companyName := c.PostForm("company_name")
	if companyName == "" {
		util.ResponseError(c, 400, "请输入公司名称")
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		util.ResponseError(c, 400, "请选择图片")
		return
	}

	adminID := c.MustGet("adminID").(uint64)
	adminIDUint := uint(adminID)

	opts := service.UploadOptions{
		ImageType:   service.ImageTypeAdminPaymentCode,
		AdminID:     adminIDUint,
		CompanyName: companyName,
	}

	result, err := service.UploadImage(file, opts)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	image := model.PaymentCodeImage{
		Path:        result.RelativePath,
		CompanyName: companyName,
		FileName:    result.FileName,
		AdminID:     adminIDUint,
		Size:        file.Size,
	}

	if err := service.CreatePaymentCodeImage(&image); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	// 将图片路径转换为完整URL
	service.ConvertPaymentCodeImageURL(&image)

	util.ResponseSuccess(c, image)
}

func GetPaymentCodeImageList(c *gin.Context) {
	page := 1
	pageSize := 10

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

	images, total, err := service.GetPaymentCodeImageList(page, pageSize)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	result := gin.H{
		"list":      images,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}

	util.ResponseSuccess(c, result)
}

func GetPaymentCodeImageDetail(c *gin.Context) {
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

	image, err := service.GetPaymentCodeImageByID(id)
	if err != nil {
		util.ResponseBizError(c, util.NewBizError(util.ErrCodePaymentCodeImageNotFound, "付款码图片不存在"))
		return
	}

	util.ResponseSuccess(c, image)
}

func UpdatePaymentCodeImage(c *gin.Context) {
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
		CompanyName string `json:"company_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	image, err := service.GetPaymentCodeImageByID(id)
	if err != nil {
		util.ResponseBizError(c, util.NewBizError(util.ErrCodePaymentCodeImageNotFound, "付款码图片不存在"))
		return
	}

	if req.CompanyName != "" {
		image.CompanyName = req.CompanyName
	}

	if err := service.UpdatePaymentCodeImage(image); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, image)
}

func DeletePaymentCodeImage(c *gin.Context) {
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

	if err := service.DeletePaymentCodeImage(id); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, nil)
}