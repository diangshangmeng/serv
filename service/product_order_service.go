package service

import (
	"time"

	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"

	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func PlaceOrder(productID uint64, userID uint64, userPhone string) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}

	if product.Status != 0 {
		return util.NewBizError(util.ErrCodeProductUnavailable, "商品不可购买")
	}

	product.Status = 1
	product.OwnerID = uint(userID)

	if err := repository.UpdateProduct(product); err != nil {
		logger.Error("更新商品状态失败", zap.Error(err), zap.Uint64("product_id", productID))
		return util.NewBizError(util.ErrCodeProductUnavailable, "商品已被其他用户购买")
	}

	_ = SetProductDetailToCache(product)
	ClearProductListCache()

	return nil
}

func CompleteProductPayment(productID uint64) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}

	if product.Status != 1 {
		return util.NewBizError(util.ErrCodeProductStatusError, "商品状态不正确")
	}

	now := time.Now()
	product.Status = 4
	product.PayTime = &now

	if err := repository.UpdateProduct(product); err != nil {
		return err
	}

	// 更新商品详情缓存，清除列表缓存
	_ = SetProductDetailToCache(product)
	ClearProductListCache()

	return nil
}

func ConfirmProductOrder(productID uint64) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}

	if product.Status != 4 {
		return util.NewBizError(util.ErrCodeProductStatusError, "商品状态不正确")
	}

	// 保存当前订单价格，用于更新 last_transaction_price
	currentPrice := product.Price

	product.Status = 0
	product.PayTime = nil
	product.LastTransactionPrice = currentPrice

	if err := repository.UpdateProduct(product); err != nil {
		logger.Error("更新商品状态失败", zap.Error(err), zap.Uint64("product_id", productID))
		return err
	}

	// 更新商品详情缓存，清除列表缓存
	_ = SetProductDetailToCache(product)
	ClearProductListCache()
	ClearMyProductsCache(product.OwnerPhone)

	return nil
}

func CancelProductOrder(productID uint64, reason string) error {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}

	if product.Status != 1 && product.Status != 2 {
		return util.NewBizError(util.ErrCodeProductStatusError, "商品状态不正确")
	}

	product.Status = 2
	product.LockReason = reason

	if err := repository.UpdateProduct(product); err != nil {
		return err
	}

	// 更新商品详情缓存，清除列表缓存
	_ = SetProductDetailToCache(product)
	ClearProductListCache()

	return nil
}

func CreateProductOrder(productID uint64, buyerID uint64, buyerPhone string) (*model.Order, *model.User, error) {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return nil, nil, util.NewBizError(util.ErrCodeProductNotFound, "商品不存在")
	}

	if product.Status != 1 {
		return nil, nil, util.NewBizError(util.ErrCodeProductUnavailable, "商品不可购买")
	}

	if uint64(product.OwnerID) == buyerID {
		return nil, nil, util.NewBizError(util.ErrCodeCannotBuyOwn, "不能购买自己的商品")
	}

	tx := config.DB.Begin()

	product.Status = 2
	if err := tx.Save(product).Error; err != nil {
		tx.Rollback()
		logger.Error("更新商品状态失败", zap.Error(err), zap.Uint64("product_id", productID))
		return nil, nil, util.NewBizError(util.ErrCodeProductUnavailable, "商品已被其他用户购买")
	}

	orderNo := util.GenerateOrderNo()
	now := time.Now()
	order := &model.Order{
		OrderNo:             orderNo,
		ProductID:           product.ID,
		ProductTitle:        product.Title,
		ProductDescription:  product.Description,
		SellerID:            product.OwnerID,
		SellerPhone:         product.OwnerPhone,
		ProductImageURL:     product.DisplayImageURL,
		PaymentCodeImageURL: product.PaymentImageURL,
		BuyerID:             uint(buyerID),
		BuyerPhone:          buyerPhone,
		Status:              0,
		Price:               product.Price,
		ExpiredAt:           now.Add(15 * time.Minute),
		OrderTime:           &now,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		logger.Error("创建订单失败", zap.Error(err), zap.String("order_no", orderNo), zap.Uint64("buyer_id", buyerID))
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("提交事务失败", zap.Error(err), zap.String("order_no", orderNo))
		return nil, nil, err
	}

	_ = SetProductDetailToCache(product)
	ClearProductListCache()
	ClearMyProductsCache(product.OwnerPhone)

	ConvertOrderImageURLs(order)
	return order, nil, nil
}
