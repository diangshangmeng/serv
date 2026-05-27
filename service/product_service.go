package service

import (
	"log"

	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func CreateProduct(product *model.Product, adminID uint, adminPhone string) error {
	product.OwnerID = adminID
	product.OwnerPhone = adminPhone
	return repository.CreateProduct(product)
}

func GetProductByID(id uint64) (*model.Product, error) {
	product, err := repository.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	ConvertProductImageURLs(product)
	return product, nil
}

func GetProductList(page int, pageSize int, cityID *uint) ([]model.Product, int64, error) {
	list, total, err := repository.GetProductList(page, pageSize, cityID)
	if err != nil {
		return nil, 0, err
	}
	ConvertProductsImageURLs(list)
	return list, total, nil
}

func GetProductListForDashboard(page int, pageSize int, cityID *uint) ([]model.Product, int64, error) {
	list, total, err := repository.GetProductListForDashboard(page, pageSize, cityID)
	if err != nil {
		return nil, 0, err
	}
	ConvertProductsImageURLs(list)
	return list, total, nil
}

func UpdateProduct(product *model.Product) error {
	return repository.UpdateProduct(product)
}

func DeleteProduct(id uint64, adminPhone string) error {
	product, err := repository.GetProductByID(id)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	if product.OwnerPhone != adminPhone {
		return util.NewBizError(util.ErrCodeProductNoPermission, "无权删除此商品")
	}
	return repository.DeleteProduct(id)
}

func UpdateProductStatus(productID uint64, status int) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	product.Status = status
	return repository.UpdateProduct(product)
}

func GetProductStatus(productID uint64) (int, error) {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return -1, util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	return product.Status, nil
}

func SetProductStatusLock(productID uint64, reason string) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	product.Status = 2
	product.LockReason = reason
	return repository.UpdateProduct(product)
}

func ResetProductToAvailable(productID uint64) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	product.Status = 0
	product.PayTime = nil
	product.LockReason = ""
	return repository.UpdateProduct(product)
}

func PublishProduct(productID uint64) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	product.Status = 1
	product.PayTime = nil
	product.LockReason = ""
	return repository.UpdateProduct(product)
}

func UnpublishProduct(productID uint64, reason string) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	product.Status = 0
	product.PayTime = nil
	product.LockReason = reason
	return repository.UpdateProduct(product)
}

func UnlockProduct(productID uint64) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}
	product.Status = 0
	product.PayTime = nil
	product.LockReason = ""
	return repository.UpdateProduct(product)
}

func GetMyProducts(ownerPhone string, page int, pageSize int, cityID *uint) ([]*model.Product, int64, error) {
	log.Printf("[GetMyProducts] 开始查询商品列表，ownerPhone=%s, page=%d, pageSize=%d", ownerPhone, page, pageSize)

	products, total, err := repository.GetProductsByOwnerPhone(ownerPhone, page, pageSize, cityID)
	if err != nil {
		log.Printf("[GetMyProducts] 查询失败，ownerPhone=%s, error=%v", ownerPhone, err)
		return nil, 0, err
	}

	for _, product := range products {
		ConvertProductImageURLs(product)
	}

	log.Printf("[GetMyProducts] 查询成功，ownerPhone=%s, total=%d", ownerPhone, total)
	return products, total, nil
}

func AppPublishProduct(productID uint64, ownerPhone string) error {
	log.Printf("[AppPublishProduct] 开始上架商品，productID=%d, ownerPhone=%s", productID, ownerPhone)

	product, err := repository.GetProductByID(productID)
	if err != nil {
		log.Printf("[AppPublishProduct] 商品不存在，productID=%d, error=%v", productID, err)
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}

	if product.OwnerPhone != ownerPhone {
		log.Printf("[AppPublishProduct] 无权操作，productID=%d, productOwner=%s, requestOwner=%s", productID, product.OwnerPhone, ownerPhone)
		return util.NewBizError(util.ErrCodeProductNoPermission, "无权操作此商品")
	}

	if product.Status != 0 {
		log.Printf("[AppPublishProduct] 状态错误，productID=%d, currentStatus=%d, expectedStatus=0", productID, product.Status)
		return util.NewBizError(util.ErrCodeProductStatusError, "商品状态不正确，只能上架未上架的商品")
	}

	product.Status = 1
	product.PayTime = nil
	product.LockReason = ""

	if err := repository.UpdateProduct(product); err != nil {
		log.Printf("[AppPublishProduct] 更新失败，productID=%d, error=%v", productID, err)
		return err
	}

	log.Printf("[AppPublishProduct] 上架成功，productID=%d", productID)
	return nil
}
