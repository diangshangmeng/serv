package controller

import (
	"voucher-platform/service"
	"voucher-platform/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UploadRequest struct {
	ImageType string `form:"image_type"`
	UserID   uint64 `form:"user_id"`
	OrderNo  string `form:"order_no"`
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		util.ResponseError(c, 400, "请选择图片")
		return
	}

	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		util.ResponseError(c, 400, "缺少必要参数：image_type")
		return
	}

	var userID uint64
	if req.UserID > 0 {
		userID = req.UserID
	} else {
		userID = c.MustGet("userID").(uint64)
	}

	opts := service.UploadOptions{
		ImageType: req.ImageType,
		UserID:    userID,
		OrderNo:   req.OrderNo,
	}

	result, err := service.UploadImage(file, opts)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	response := gin.H{
		"url":  result.FullURL,
		"type": req.ImageType,
	}

	util.ResponseSuccess(c, response)
}

func UploadAuthImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		util.ResponseError(c, 400, "请选择图片")
		return
	}

	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		util.ResponseError(c, 400, "缺少必要参数：image_type")
		return
	}

	userID := c.MustGet("userID").(uint64)

	opts := service.UploadOptions{
		ImageType: req.ImageType,
		UserID:    userID,
	}

	result, err := service.UploadImage(file, opts)
	if err != nil {
		util.ResponseBizError(c, err)
		return
	}

	response := gin.H{
		"url":  result.FullURL,
		"type": req.ImageType,
	}

	util.ResponseSuccess(c, response)
}

func UploadOrderImage(c *gin.Context) {
	util.GetLogger().Info("UploadOrderImage - 开始处理请求")

	file, err := c.FormFile("file")
	if err != nil {
		util.GetLogger().Error("UploadOrderImage - 获取文件失败", zap.Error(err))
		util.ResponseError(c, 400, "请选择图片")
		return
	}

	util.GetLogger().Info("UploadOrderImage - 获取文件成功",
		util.StringField("file_name", file.Filename),
		util.Int64Field("file_size", file.Size))

	var req UploadRequest
	if err := c.ShouldBind(&req); err != nil {
		util.GetLogger().Error("UploadOrderImage - 绑定参数失败", zap.Error(err))
		util.ResponseError(c, 400, "缺少必要参数")
		return
	}

	if req.ImageType == "" || req.ImageType != "payment_voucher" {
		req.ImageType = "payment_voucher"
	}

	if req.OrderNo == "" {
		util.GetLogger().Error("UploadOrderImage - 缺少订单号")
		util.ResponseError(c, 400, "缺少订单号")
		return
	}

	userID := c.MustGet("userID").(uint64)
	util.GetLogger().Info("UploadOrderImage - 从JWT获取用户ID", util.Uint64Field("user_id", userID))

	opts := service.UploadOptions{
		ImageType: req.ImageType,
		UserID:    userID,
		OrderNo:   req.OrderNo,
	}

	result, err := service.UploadImage(file, opts)
	if err != nil {
		util.GetLogger().Error("UploadOrderImage - 上传失败", zap.Error(err))
		util.ResponseBizError(c, err)
		return
	}

	util.GetLogger().Info("UploadOrderImage - 上传成功", util.StringField("url", result.FullURL))

	response := gin.H{
		"url":      result.FullURL,
		"type":     req.ImageType,
		"order_no": req.OrderNo,
	}

	util.ResponseSuccess(c, response)
}