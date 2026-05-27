package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/util"
)

// ImageType 图片类型枚举
const (
	ImageTypeIDCardFront     = "id_card_front"
	ImageTypeIDCardBack      = "id_card_back"
	ImageTypeBusinessLicense = "business_license"
	ImageTypePaymentCode     = "payment_code"
	ImageTypePaymentVoucher  = "payment_voucher"
	ImageTypeProductImage    = "product_image"
	ImageTypeAdminPaymentCode = "admin_payment_code"
)

// UploadOptions 上传选项
type UploadOptions struct {
	ImageType string
	UserID    uint64
	AdminID   uint
	OrderNo   string
	CompanyName string
}

// UploadResult 上传结果
type UploadResult struct {
	RelativePath string
	FullURL      string
	FileName     string
	Size         int64
}

// UploadImage 统一图片上传处理
func UploadImage(file *multipart.FileHeader, opts UploadOptions) (*UploadResult, error) {
	if file == nil {
		return nil, util.NewBizError(util.ErrCodeParamInvalid, "请选择图片")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !util.ValidateImageFormat(ext) {
		return nil, util.NewBizError(util.ErrCodeParamInvalid, "图片格式不支持")
	}

	if !util.ValidateImageSize(file.Size, config.AppConfig.MaxImageSize) {
		return nil, util.NewBizError(util.ErrCodeParamInvalid, "图片大小不能超过2MB")
	}

	uploadPath, fileName, err := getUploadPathAndName(opts, ext)
	if err != nil {
		return nil, err
	}

	relativePath, err := util.SaveImage(file, uploadPath, fileName)
	if err != nil {
		return nil, util.NewBizError(util.ErrCodeSystem, "上传失败")
	}

	return &UploadResult{
		RelativePath: relativePath,
		FullURL:      util.GetFullImageURL(relativePath),
		FileName:     fileName,
		Size:         file.Size,
	}, nil
}

// DeleteImage 统一图片删除处理
func DeleteImage(relativePath string) error {
	if relativePath == "" {
		return nil
	}
	return util.DeleteImage(relativePath)
}

// BatchUploadProductImages 批量上传商品图片
func BatchUploadProductImages(files []*multipart.FileHeader, adminID uint) ([]*UploadResult, []string, error) {
	var results []*UploadResult
	var errors []string

	for _, file := range files {
		opts := UploadOptions{
			ImageType: ImageTypeProductImage,
			AdminID:   adminID,
		}
		result, err := UploadImage(file, opts)
		if err != nil {
			errors = append(errors, fmt.Sprintf("图片 %s: %s", file.Filename, err.Error()))
			continue
		}
		results = append(results, result)
	}

	return results, errors, nil
}

// ConvertUserImageURLs 统一转换用户图片URL
func ConvertUserImageURLs(user *model.User) {
	if user == nil {
		return
	}
	user.IDCardFront = util.GetFullImageURL(user.IDCardFront)
	user.IDCardBack = util.GetFullImageURL(user.IDCardBack)
	user.BusinessLicense = util.GetFullImageURL(user.BusinessLicense)
	user.PaymentCode = util.GetFullImageURL(user.PaymentCode)
}

// ConvertUsersImageURLs 批量转换用户图片URL
func ConvertUsersImageURLs(users []model.User) {
	for i := range users {
		ConvertUserImageURLs(&users[i])
	}
}

// ConvertProductImageURLs 转换商品图片URL
func ConvertProductImageURLs(product *model.Product) {
	if product == nil {
		return
	}
	product.DisplayImageURL = util.GetFullImageURL(product.DisplayImageURL)
	product.PaymentImageURL = util.GetFullImageURL(product.PaymentImageURL)
}

// ConvertProductsImageURLs 批量转换商品图片URL
func ConvertProductsImageURLs(products []model.Product) {
	for i := range products {
		ConvertProductImageURLs(&products[i])
	}
}

// ConvertProductImageListURLs 转换商品图片列表URL
func ConvertProductImageListURLs(images []model.ProductImage) {
	for i := range images {
		images[i].Path = util.GetFullImageURL(images[i].Path)
	}
}

// ConvertProductImageURL 转换单个商品图片URL
func ConvertProductImageURL(image *model.ProductImage) {
	if image == nil {
		return
	}
	image.Path = util.GetFullImageURL(image.Path)
}

// ConvertPaymentCodeImageURLs 转换支付码图片列表URL
func ConvertPaymentCodeImageURLs(images []model.PaymentCodeImage) {
	for i := range images {
		images[i].Path = util.GetFullImageURL(images[i].Path)
	}
}

// ConvertPaymentCodeImageURL 转换单个支付码图片URL
func ConvertPaymentCodeImageURL(image *model.PaymentCodeImage) {
	if image == nil {
		return
	}
	image.Path = util.GetFullImageURL(image.Path)
}

// ConvertOrderImageURLs 转换订单图片URL
func ConvertOrderImageURLs(order *model.Order) {
	if order == nil {
		return
	}
	order.ProductImageURL = util.GetFullImageURL(order.ProductImageURL)
	order.PaymentCodeImageURL = util.GetFullImageURL(order.PaymentCodeImageURL)
	order.PaymentVoucher = util.GetFullImageURL(order.PaymentVoucher)
}

// ConvertOrdersImageURLs 批量转换订单图片URL
func ConvertOrdersImageURLs(orders []model.Order) {
	for i := range orders {
		ConvertOrderImageURLs(&orders[i])
	}
}

// getUploadPathAndName 获取上传路径和文件名
func getUploadPathAndName(opts UploadOptions, ext string) (string, string, error) {
	timestamp := time.Now().Format("20060102150405")

	switch opts.ImageType {
	case ImageTypeIDCardFront, ImageTypeIDCardBack, ImageTypeBusinessLicense, ImageTypePaymentCode:
		if opts.UserID == 0 {
			return "", "", util.NewBizError(util.ErrCodeParamInvalid, "缺少用户ID")
		}
		uploadPath := filepath.Join(config.AppConfig.ImageUploadPath, fmt.Sprintf("%d", opts.UserID), "auth")
		fileName := fmt.Sprintf("%s%s", opts.ImageType, ext)
		return uploadPath, fileName, nil

	case ImageTypePaymentVoucher:
		if opts.OrderNo == "" {
			return "", "", util.NewBizError(util.ErrCodeParamInvalid, "缺少订单号")
		}
		uploadPath := filepath.Join(config.AppConfig.ImageUploadPath, opts.OrderNo, "order")
		fileName := fmt.Sprintf("%s_%s%s", opts.OrderNo, timestamp, ext)
		return uploadPath, fileName, nil

	case ImageTypeProductImage:
		uploadPath := filepath.Join(config.AppConfig.ImageUploadPath, "dianshangmeng", "products")
		randStr := fmt.Sprintf("%d", time.Now().UnixNano()%10000)
		fileName := fmt.Sprintf("%s_%s%s", timestamp, randStr, ext)
		return uploadPath, fileName, nil

	case ImageTypeAdminPaymentCode:
		uploadPath := filepath.Join(config.AppConfig.ImageUploadPath, "dianshangmeng", "auth")
		// 如果有公司名，就用公司名的安全版本作为文件名的一部分
		fileName := fmt.Sprintf("%s%s", timestamp, ext)
		if opts.CompanyName != "" {
			// 简单清理公司名，只保留字母数字，避免文件名问题
			cleanName := ""
			for _, r := range opts.CompanyName {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
					cleanName += string(r)
				}
			}
			if cleanName != "" {
				fileName = fmt.Sprintf("%s_%s%s", timestamp, cleanName, ext)
			}
		}
		return uploadPath, fileName, nil

	default:
		return "", "", util.NewBizError(util.ErrCodeParamInvalid, "不支持的图片类型")
	}
}
