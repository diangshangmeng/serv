package controller

import (
	"mime/multipart"
	"strconv"
	"sync"

	"voucher-platform/model"
	"voucher-platform/service"
	"voucher-platform/util"

	"github.com/gin-gonic/gin"
)

func BatchUploadProductImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		util.ResponseError(c, 400, "请选择图片")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		util.ResponseError(c, 400, "请选择图片")
		return
	}

	if len(files) > 20 {
		util.ResponseError(c, 400, "一次最多上传20张图片")
		return
	}

	adminID := c.MustGet("adminID").(uint64)
	adminIDUint := uint(adminID)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var images []model.ProductImage
	var errs []string

	for _, file := range files {
		wg.Add(1)
		go func(f *multipart.FileHeader) {
			defer wg.Done()

			opts := service.UploadOptions{
				ImageType: service.ImageTypeProductImage,
				AdminID:   adminIDUint,
			}

			result, err := service.UploadImage(f, opts)
			if err != nil {
				mu.Lock()
				errs = append(errs, err.Error())
				mu.Unlock()
				return
			}

			image := model.ProductImage{
				Path:     result.RelativePath,
				Name:     f.Filename,
				FileName: result.FileName,
				Tags:     "",
				IsUsed:   false,
				AdminID:  adminIDUint,
				Size:     f.Size,
			}

			mu.Lock()
			images = append(images, image)
			mu.Unlock()
		}(file)
	}

	wg.Wait()

	if len(images) > 0 {
		for i := range images {
			if err := service.CreateProductImage(&images[i]); err != nil {
				mu.Lock()
				errs = append(errs, util.NewBizError(util.ErrCodeImageNotFound, "图片保存失败").Error())
				mu.Unlock()
			}
			// 将图片路径转换为完整URL
			service.ConvertProductImageURL(&images[i])
		}
	}

	result := gin.H{
		"success_count": len(images),
		"failed_count":  len(errs),
		"images":        images,
		"errors":        errs,
	}

	util.ResponseSuccess(c, result)
}

func GetProductImageList(c *gin.Context) {
	page := 1
	pageSize := 500

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

	var isUsed *bool
	if s := c.Query("is_used"); s != "" {
		if parsed, err := strconv.ParseBool(s); err == nil {
			isUsed = &parsed
		}
	}

	images, total, err := service.GetProductImageList(page, pageSize, isUsed)
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

func GetProductImageDetail(c *gin.Context) {
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

	image, err := service.GetProductImageByID(id)
	if err != nil {
		util.ResponseBizError(c, util.NewBizError(util.ErrCodeImageNotFound, "图片不存在"))
		return
	}

	util.ResponseSuccess(c, image)
}

func UpdateProductImage(c *gin.Context) {
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
		Name string `json:"name"`
		Tags string `json:"tags"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	image, err := service.GetProductImageByID(id)
	if err != nil {
		util.ResponseBizError(c, util.NewBizError(util.ErrCodeImageNotFound, "图片不存在"))
		return
	}

	if req.Name != "" {
		image.Name = req.Name
	}
	if req.Tags != "" {
		image.Tags = req.Tags
	}

	if err := service.UpdateProductImage(image); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, image)
}

func MarkProductImageAsUsed(c *gin.Context) {
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

	if err := service.MarkProductImageAsUsed(id); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, nil)
}

func DeleteProductImage(c *gin.Context) {
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

	if err := service.DeleteProductImage(id); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, nil)
}

func BatchDeleteProductImages(c *gin.Context) {
	var req struct {
		IDs []uint64 `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.ResponseError(c, 400, "参数错误")
		return
	}

	if len(req.IDs) == 0 {
		util.ResponseError(c, 400, "请选择要删除的图片")
		return
	}

	if err := service.BatchDeleteProductImages(req.IDs); err != nil {
		util.ResponseBizError(c, err)
		return
	}

	util.ResponseSuccess(c, nil)
}