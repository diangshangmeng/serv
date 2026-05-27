package service

import (
	"log"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func CreateProductImage(image *model.ProductImage) error {
	return repository.CreateProductImage(image)
}

func BatchCreateProductImages(images []*model.ProductImage) error {
	return repository.BatchCreateProductImages(images)
}

func GetProductImageByID(id uint64) (*model.ProductImage, error) {
	image, err := repository.GetProductImageByID(id)
	if err != nil {
		return nil, err
	}
	ConvertProductImageURL(image)
	return image, nil
}

func GetProductImageList(page int, pageSize int, isUsed *bool) ([]model.ProductImage, int64, error) {
	images, total, err := repository.GetProductImageList(page, pageSize, isUsed)
	if err != nil {
		return nil, 0, err
	}
	ConvertProductImageListURLs(images)
	return images, total, nil
}

func UpdateProductImage(image *model.ProductImage) error {
	return repository.UpdateProductImage(image)
}

func MarkProductImageAsUsed(id uint64) error {
	image, err := repository.GetProductImageByID(id)
	if err != nil {
		return util.NewBizError(util.ErrCodeImageNotFound, "图片不存在")
	}

	if image.IsUsed {
		return util.NewBizError(util.ErrCodeImageAlreadyUsed, "图片已被使用")
	}

	image.IsUsed = true
	return repository.UpdateProductImage(image)
}

func DeleteProductImage(id uint64) error {
	image, err := repository.GetProductImageByID(id)
	if err != nil {
		return util.NewBizError(util.ErrCodeImageNotFound, "图片不存在")
	}

	if err := DeleteImage(image.Path); err != nil {
		log.Printf("[DeleteProductImage] 删除物理文件失败，id=%d, path=%s, error=%v", id, image.Path, err)
	}

	return repository.DeleteProductImage(id)
}

func BatchDeleteProductImages(ids []uint64) error {
	for _, id := range ids {
		image, err := repository.GetProductImageByID(id)
		if err != nil {
			log.Printf("[BatchDeleteProductImages] 获取图片失败，id=%d, error=%v", id, err)
			continue
		}

		if err := DeleteImage(image.Path); err != nil {
			log.Printf("[BatchDeleteProductImages] 删除物理文件失败，id=%d, path=%s, error=%v", id, image.Path, err)
		}
	}

	return repository.BatchDeleteProductImages(ids)
}
