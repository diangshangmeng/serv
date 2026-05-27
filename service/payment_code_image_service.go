package service

import (
	"log"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func CreatePaymentCodeImage(image *model.PaymentCodeImage) error {
	return repository.CreatePaymentCodeImage(image)
}

func GetPaymentCodeImageByID(id uint64) (*model.PaymentCodeImage, error) {
	image, err := repository.GetPaymentCodeImageByID(id)
	if err != nil {
		return nil, util.NewBizError(util.ErrCodePaymentCodeImageNotFound, "付款码图片不存在")
	}
	ConvertPaymentCodeImageURL(image)
	return image, nil
}

func GetPaymentCodeImageList(page int, pageSize int) ([]model.PaymentCodeImage, int64, error) {
	images, total, err := repository.GetPaymentCodeImageList(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertPaymentCodeImageURLs(images)
	return images, total, nil
}

func UpdatePaymentCodeImage(image *model.PaymentCodeImage) error {
	_, err := repository.GetPaymentCodeImageByID(uint64(image.ID))
	if err != nil {
		return util.NewBizError(util.ErrCodePaymentCodeImageNotFound, "付款码图片不存在")
	}
	return repository.UpdatePaymentCodeImage(image)
}

func DeletePaymentCodeImage(id uint64) error {
	image, err := repository.GetPaymentCodeImageByID(id)
	if err != nil {
		return util.NewBizError(util.ErrCodePaymentCodeImageNotFound, "付款码图片不存在")
	}

	if err := DeleteImage(image.Path); err != nil {
		log.Printf("[DeletePaymentCodeImage] 删除物理文件失败，id=%d, path=%s, error=%v", id, image.Path, err)
	}

	return repository.DeletePaymentCodeImage(id)
}
