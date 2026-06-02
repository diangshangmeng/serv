package service

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

const (
	ORDER_STATUS_UNPAID    = 0
	ORDER_STATUS_PAID      = 1
	ORDER_STATUS_COMPLETED = 2
	ORDER_STATUS_CANCELLED = 3
	ORDER_STATUS_TIMEOUT   = 4
)

func UploadPaymentVoucher(orderID uint, voucherURL string) error {
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != ORDER_STATUS_UNPAID {
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确，只能为未付款状态")
	}

	order.PaymentVoucher = voucherURL
	order.Status = ORDER_STATUS_PAID

	return repository.UpdateOrder(order)
}

func ConfirmOrderByID(orderID, sellerID uint) error {
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	if order.SellerID != sellerID {
		return util.NewBizError(util.ErrCodeNoPermission, "无权操作此订单")
	}

	if order.Status != ORDER_STATUS_PAID {
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确，只能为已付款待确认状态")
	}

	return repository.UpdateOrderStatus(orderID, ORDER_STATUS_COMPLETED)
}

func CancelOrderByID(orderID, userID uint, isSeller bool) error {
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	if isSeller {
		if order.SellerID != userID {
			return util.NewBizError(util.ErrCodeNoPermission, "无权操作此订单")
		}
		if order.Status != ORDER_STATUS_UNPAID && order.Status != ORDER_STATUS_PAID {
			return util.NewBizError(util.ErrCodeOrderStatusError, "卖家只能取消未付款或已付款待确认的订单")
		}
	} else {
		if order.BuyerID != userID {
			return util.NewBizError(util.ErrCodeNoPermission, "无权操作此订单")
		}
		if order.Status != ORDER_STATUS_UNPAID {
			return util.NewBizError(util.ErrCodeOrderStatusError, "买家只能取消未付款的订单")
		}
	}

	return repository.UpdateOrderStatus(orderID, ORDER_STATUS_CANCELLED)
}

func CloseExpiredOrders() error {
	orders, err := repository.GetExpiredOrders()
	if err != nil {
		return err
	}

	for _, order := range orders {
		tx := config.DB.Begin()

		err := repository.UpdateOrderStatus(order.ID, ORDER_STATUS_TIMEOUT)
		if err != nil {
			tx.Rollback()
			return err
		}

		product, err := repository.GetProductByID(uint64(order.ProductID))
		if err != nil {
			tx.Rollback()
			return err
		}
		fmt.Println("bbbbb")
		product.Status = 1
		err = tx.Save(product).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		
		// 更新商品详情缓存，清除列表缓存
		_ = SetProductDetailToCache(product)
		ClearProductListCache()
		// 清除卖家和买家的"我的商品"缓存
		ClearMyProductsCache(product.OwnerPhone)
		ClearMyProductsCache(order.BuyerPhone)
	}

	return nil
}

func GetOrderDetailByID(orderID uint) (*model.Order, error) {
	order, err := repository.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}
	ConvertOrderImageURLs(order)
	return order, nil
}

func GetBuyerOrders(phone string, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetOrdersByBuyerID(phone, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func GetSellerOrders(phone string, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetOrdersBySellerPhone(phone, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func CreateOrder(buyerID uint64, listingID uint64) (*model.Order, *model.User, error) {
	err := CheckUserAuthStatus(buyerID)
	if err != nil {
		return nil, nil, err
	}

	listing, err := repository.GetListingByID(listingID)
	if err != nil {
		return nil, nil, err
	}

	if listing.Status != 0 {
		return nil, nil, util.NewBizError(util.ErrCodeListingStatusError, "挂单已被购买")
	}

	coupon, err := repository.GetCouponByID(uint64(listing.CouponID))
	if err != nil {
		return nil, nil, err
	}

	seller, err := repository.GetUserByID(uint64(coupon.UserID))
	if err != nil {
		return nil, nil, err
	}

	buyer, err := repository.GetUserByID(buyerID)
	if err != nil {
		return nil, nil, err
	}

	if buyer.CityID != seller.CityID {
		return nil, nil, util.NewBizError(util.ErrCodeNotSameCity, "只能购买同城市的券码")
	}

	if buyer.ID == seller.ID {
		return nil, nil, util.NewBizError(util.ErrCodeCannotBuyOwn, "不能购买自己的券码")
	}

	tx := config.DB.Begin()

	listing.Status = 1
	err = tx.Save(listing).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	coupon.Status = 2
	err = tx.Save(coupon).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	product, err := repository.GetProductByID(uint64(listing.ProductID))
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	orderNo := util.GenerateOrderNo()
	now := time.Now()
	order := &model.Order{
		OrderNo:             orderNo,
		ProductID:           product.ID,
		ProductTitle:        product.Title,
		ProductDescription:  product.Description,
		SellerID:            seller.ID,
		SellerPhone:         seller.Phone,
		ProductImageURL:     product.DisplayImageURL,
		PaymentCodeImageURL: product.PaymentImageURL,
		BuyerID:             uint(buyerID),
		BuyerPhone:          buyer.Phone,
		Status:              ORDER_STATUS_UNPAID,
		Price:               int64(listing.SellingPrice * 100),
		ExpiredAt:           now.Add(15 * time.Minute),
		OrderTime:           &now,
	}

	err = tx.Create(order).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	tx.Commit()

	ctx := context.Background()
	redisKey := fmt.Sprintf("order_lock:%s", orderNo)
	config.RedisClient.Set(ctx, redisKey, "locked", 15*time.Minute)

	ConvertOrderImageURLs(order)
	return order, seller, nil
}

func UploadVoucher(orderNo string, paymentVoucher string, buyerID uint64) error {
	util.GetLogger().Info("UploadVoucher - 开始上传凭证",
		util.StringField("order_no", orderNo),
		util.Uint64Field("buyer_id", buyerID))

	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		util.GetLogger().Error("UploadVoucher - 查询订单失败", zap.Error(err),
			util.StringField("order_no", orderNo))
		return err
	}
	util.GetLogger().Info("UploadVoucher - 查询订单成功",
		util.StringField("order_no", order.OrderNo),
		util.Uint64Field("buyer_id", uint64(order.BuyerID)),
		util.IntField("status", order.Status))

	if uint64(order.BuyerID) != buyerID {
		util.GetLogger().Error("UploadVoucher - 权限验证失败",
			util.StringField("order_no", orderNo),
			util.Uint64Field("order_buyer_id", uint64(order.BuyerID)),
			util.Uint64Field("request_buyer_id", buyerID))
		return util.NewBizError(util.ErrCodeNoPermission, "无权操作此订单")
	}
	util.GetLogger().Info("UploadVoucher - 权限验证成功",
		util.StringField("order_no", orderNo),
		util.Uint64Field("buyer_id", buyerID))

	if order.Status != ORDER_STATUS_UNPAID {
		util.GetLogger().Error("UploadVoucher - 订单状态验证失败",
			util.StringField("order_no", orderNo),
			util.IntField("current_status", order.Status),
			util.IntField("expected_status", ORDER_STATUS_UNPAID))
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确")
	}
	util.GetLogger().Info("UploadVoucher - 订单状态验证成功",
		util.StringField("order_no", orderNo),
		util.IntField("status", order.Status))

	order.PaymentVoucher = paymentVoucher
	order.Status = ORDER_STATUS_PAID

	err = repository.UpdateOrder(order)
	if err != nil {
		util.GetLogger().Error("UploadVoucher - 更新订单失败", zap.Error(err),
			util.StringField("order_no", orderNo))
		return err
	}

	util.GetLogger().Info("UploadVoucher - 更新订单成功",
		util.StringField("order_no", orderNo),
		util.IntField("new_status", ORDER_STATUS_PAID))
	return nil
}

func ConfirmOrder(orderNo string, sellerPhone string) error {
	util.GetLogger().Info("ConfirmOrder - 开始确认订单",
		util.StringField("order_no", orderNo),
		util.StringField("seller_phone", sellerPhone))

	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		util.GetLogger().Error("ConfirmOrder - 查询订单失败", zap.Error(err),
			util.StringField("order_no", orderNo))
		return err
	}

	if order.SellerPhone != sellerPhone {
		util.GetLogger().Error("ConfirmOrder - 权限验证失败",
			util.StringField("order_no", orderNo),
			util.StringField("order_seller_phone", order.SellerPhone),
			util.StringField("request_seller_phone", sellerPhone))
		return util.NewBizError(util.ErrCodeNoPermission, "无权操作此订单")
	}

	if order.Status != ORDER_STATUS_PAID {
		util.GetLogger().Error("ConfirmOrder - 订单状态验证失败",
			util.StringField("order_no", orderNo),
			util.IntField("current_status", order.Status),
			util.IntField("expected_status", ORDER_STATUS_PAID))
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确")
	}

	tx := config.DB.Begin()

	// 尝试查找 listing，如果找不到也继续执行
	listing, err := repository.GetListingByProductID(uint64(order.ProductID))
	if err == nil {
		// 如果找到 listing，尝试更新 coupon
		coupon, err := repository.GetCouponByID(uint64(listing.CouponID))
		if err == nil {
			coupon.UserID = uint(order.BuyerID)
			coupon.Status = 0
			err = tx.Save(coupon).Error
			if err != nil {
				util.GetLogger().Warn("ConfirmOrder - 更新 coupon 失败，继续执行", zap.Error(err),
					util.StringField("order_no", orderNo))
			}
		} else {
			util.GetLogger().Warn("ConfirmOrder - 查询 coupon 失败，继续执行", zap.Error(err),
				util.StringField("order_no", orderNo))
		}
	} else {
		util.GetLogger().Warn("ConfirmOrder - 查询 listing 失败，继续执行", zap.Error(err),
			util.StringField("order_no", orderNo))
	}

	product, err := repository.GetProductByID(uint64(order.ProductID))
	if err != nil {
		tx.Rollback()
		return err
	}
	// 记录原卖家手机号，用于清除其"我的商品"缓存
	originalOwnerPhone := product.OwnerPhone
	product.Status = 0
	product.OwnerID = order.BuyerID
	product.OwnerPhone = order.BuyerPhone
	
	buyer, err := repository.GetUserByID(uint64(order.BuyerID))
	if err != nil {
		tx.Rollback()
		return err
	}
	product.PaymentImageURL = buyer.PaymentCode
	
	util.GetLogger().Info("ConfirmOrder - 更新产品所有者",
		util.StringField("order_no", orderNo),
		util.Uint64Field("product_id", uint64(product.ID)),
		util.Uint64Field("new_owner_id", uint64(product.OwnerID)),
		util.StringField("new_owner_phone", product.OwnerPhone))
	err = tx.Save(product).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	order.Status = ORDER_STATUS_COMPLETED
	now := time.Now()
	order.ExpiredAt = now
	err = tx.Save(order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("order_lock:%s", orderNo)
	config.RedisClient.Del(ctx, redisKey)

	tx.Commit()
	
	// 更新商品详情缓存，清除列表缓存
	_ = SetProductDetailToCache(product)
	ClearProductListCache()
	// 清除原卖家的"我的商品"缓存（商品已不属于该卖家）
	ClearMyProductsCache(originalOwnerPhone)
	// 清除新买家的"我的商品"缓存（商品现在属于该买家）
	ClearMyProductsCache(order.BuyerPhone)
	
	util.GetLogger().Info("ConfirmOrder - 确认订单成功",
		util.StringField("order_no", orderNo),
		util.StringField("seller_phone", sellerPhone))
	return nil
}

func SellerCancelOrder(orderNo string, sellerPhone string) error {
	util.GetLogger().Info("SellerCancelOrder - 开始卖家取消订单",
		util.StringField("order_no", orderNo),
		util.StringField("seller_phone", sellerPhone))

	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		util.GetLogger().Error("SellerCancelOrder - 查询订单失败", zap.Error(err),
			util.StringField("order_no", orderNo))
		return err
	}
	util.GetLogger().Info("SellerCancelOrder - 查询订单成功",
		util.StringField("order_no", order.OrderNo),
		util.StringField("seller_phone", order.SellerPhone),
		util.IntField("status", order.Status))

	if order.SellerPhone != sellerPhone {
		util.GetLogger().Error("SellerCancelOrder - 权限验证失败",
			util.StringField("order_no", orderNo),
			util.StringField("order_seller_phone", order.SellerPhone),
			util.StringField("request_seller_phone", sellerPhone))
		return util.NewBizError(util.ErrCodeNoPermission, "无权操作此订单")
	}
	util.GetLogger().Info("SellerCancelOrder - 权限验证成功",
		util.StringField("order_no", orderNo),
		util.StringField("seller_phone", sellerPhone))

	if order.Status != ORDER_STATUS_PAID {
		util.GetLogger().Error("SellerCancelOrder - 订单状态验证失败",
			util.StringField("order_no", orderNo),
			util.IntField("current_status", order.Status),
			util.IntField("expected_status", ORDER_STATUS_PAID))
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确")
	}
	util.GetLogger().Info("SellerCancelOrder - 订单状态验证成功",
		util.StringField("order_no", orderNo),
		util.IntField("status", order.Status))

	tx := config.DB.Begin()

	now := time.Now()
	order.Status = ORDER_STATUS_UNPAID
	order.ExpiredAt = now.Add(15 * time.Minute)
	order.PaymentVoucher = ""
	err = tx.Save(order).Error
	if err != nil {
		util.GetLogger().Error("SellerCancelOrder - 更新订单失败", zap.Error(err),
			util.StringField("order_no", orderNo))
		tx.Rollback()
		return err
	}
	util.GetLogger().Info("SellerCancelOrder - 更新订单成功",
		util.StringField("order_no", orderNo),
		util.IntField("new_status", ORDER_STATUS_UNPAID),
		zap.Time("expired_at", order.ExpiredAt))

	tx.Commit()
	
	// 清除卖家和买家的"我的商品"缓存
	ClearMyProductsCache(order.SellerPhone)
	ClearMyProductsCache(order.BuyerPhone)
	// 清除商品列表缓存
	ClearProductListCache()
	
	util.GetLogger().Info("SellerCancelOrder - 卖家取消订单成功",
		util.StringField("order_no", orderNo),
		util.StringField("seller_phone", sellerPhone))
	return nil
}

func CancelOrder(orderNo string) error {
	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		return err
	}

	if order.Status != ORDER_STATUS_UNPAID && order.Status != ORDER_STATUS_PAID {
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不允许取消")
	}

	tx := config.DB.Begin()

	listing, err := repository.GetListingByProductID(uint64(order.ProductID))
	if err != nil {
		tx.Rollback()
		return err
	}

	coupon, err := repository.GetCouponByID(uint64(listing.CouponID))
	if err != nil {
		tx.Rollback()
		return err
	}

	coupon.Status = 0
	err = tx.Save(coupon).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	listing.Status = 0
	err = tx.Save(listing).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	product, err := repository.GetProductByID(uint64(order.ProductID))
	if err != nil {
		tx.Rollback()
		return err
	}
	product.Status = 1
	err = tx.Save(product).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	order.Status = ORDER_STATUS_CANCELLED
	err = tx.Save(order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	ctx := context.Background()
	redisKey := fmt.Sprintf("order_lock:%s", orderNo)
	config.RedisClient.Del(ctx, redisKey)

	tx.Commit()
	
	// 更新商品详情缓存，清除列表缓存
	_ = SetProductDetailToCache(product)
	ClearProductListCache()
	// 清除卖家和买家的"我的商品"缓存
	ClearMyProductsCache(product.OwnerPhone)
	ClearMyProductsCache(order.BuyerPhone)
	return nil
}

func GetMyBuyOrders(phone string, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetOrdersByBuyerID(phone, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func GetMySellOrders(phone string, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetOrdersBySellerPhone(phone, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func GetOrderDetail(orderNo string, userID uint64) (*model.Order, error) {
	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}

	if uint64(order.BuyerID) != userID && uint64(order.SellerID) != userID {
		return nil, util.NewBizError(util.ErrCodeNoPermission, "无权查看此订单")
	}

	ConvertOrderImageURLs(order)
	return order, nil
}

func GetBuyerOrdersByPhone(phone string, status *int, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetBuyerOrdersByPhone(phone, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func GetSellerOrdersByPhone(phone string, status *int, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetSellerOrdersByPhone(phone, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func GetAllOrdersByPhone(phone string, status *int, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.GetAllOrdersByPhone(phone, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func AdminGetOrders(orderNo, phone string, status *int, startTime, endTime *time.Time, page, pageSize int) ([]model.Order, int64, error) {
	orders, total, err := repository.AdminGetOrders(orderNo, phone, status, startTime, endTime, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertOrdersImageURLs(orders)
	return orders, total, nil
}

func GetOrderByOrderNo(orderNo string) (*model.Order, error) {
	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		return nil, err
	}
	ConvertOrderImageURLs(order)
	return order, nil
}

func GetOrderStats() (map[string]int64, error) {
	return repository.GetOrderStats()
}

func GetActiveOrderByProductID(productID uint64) (*model.Order, error) {
	order, err := repository.GetActiveOrderByProductID(productID)
	if err != nil {
		return nil, err
	}
	if order != nil {
		ConvertOrderImageURLs(order)
	}
	return order, nil
}
