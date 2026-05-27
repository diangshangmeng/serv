package service

import (
	"context"
	"strings"
	"time"

	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"

	"go.uber.org/zap"
)

func StartOrderTimeoutChecker() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			checkTimeoutOrders()
			checkProductTimeout()
		}
	}()
}

func checkProductTimeout() {
	ctx := context.Background()
	now := time.Now()

	var products []model.Product
	config.DB.Where("status = ?", 4).Find(&products)

	for _, product := range products {
		if product.Status == 4 && product.PayTime != nil {
			if now.Sub(*product.PayTime) >= 30*time.Minute {
				util.GetLogger().Warn("product_confirm_timeout",
					util.Uint64Field("product_id", uint64(product.ID)),
					util.StringField("status", "4->2"),
				)
				SetProductStatusLock(uint64(product.ID), "确认超时")
			}
		}
	}

	_ = ctx
}

func checkTimeoutOrders() {
	var orders []model.Order
	orders_status := 0 // 0:待支付 1:已支付待确认 2:已完成 3:已取消 4:已超时
	now := time.Now()
	config.DB.Where("status = ? AND expired_at <= ?", orders_status, now).Find(&orders)

	util.GetLogger().Info("check_timeout_orders_start",
		util.IntField("found_orders", len(orders)),
		zap.Time("current_time", now),
	)

	for _, order := range orders {
		util.GetLogger().Info("order_timeout_closed",
			util.StringField("order_no", order.OrderNo),
			util.Uint64Field("buyer_id", uint64(order.BuyerID)),
			zap.Time("expired_at", order.ExpiredAt),
			zap.Time("created_at", order.CreatedAt),
			zap.Time("now", now),
		)
		config.DB.Model(&order).Update("status", 4)
		err := changeProductStatusToAvailable(uint64(order.ProductID))
		if err != nil {
			util.GetLogger().Error("change_product_status_failed",
				util.StringField("order_no", order.OrderNo),
				util.Uint64Field("product_id", uint64(order.ProductID)),
				zap.Error(err),
			)
		} else {
			util.GetLogger().Info("change_product_status_success",
				util.StringField("order_no", order.OrderNo),
				util.Uint64Field("product_id", uint64(order.ProductID)),
			)
		}
	}

	ctx := context.Background()
	keys, err := config.RedisClient.Keys(ctx, "order_lock:*").Result()
	if err != nil {
		util.GetLogger().Error("get_order_lock_keys_failed",
			zap.Error(err),
		)
		return
	}

	for _, key := range keys {
		ttl, err := config.RedisClient.TTL(ctx, key).Result()
		if err != nil {
			util.GetLogger().Error("get_ttl_failed",
				util.StringField("key", key),
				zap.Error(err),
			)
			continue
		}

		if ttl < 0 {
			orderNo := strings.TrimPrefix(key, "order_lock:")
			CancelOrder(orderNo)
		}
	}

}

func changeProductStatusToAvailable(productID uint64) error {
	tx := config.DB.Begin()

	product, err := repository.GetProductByID(productID)
	if err != nil {
		util.GetLogger().Error("get_product_failed",
			util.Uint64Field("product_id", productID),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}
	product.Status = 1
	err = tx.Save(product).Error
	if err != nil {
		util.GetLogger().Error("save_product_failed",
			util.Uint64Field("product_id", productID),
			zap.Error(err),
		)
		tx.Rollback()
		return err
	}

	tx.Commit()
	util.GetLogger().Info("product_status_changed_to_available",
		util.Uint64Field("product_id", productID),
	)
	return nil
}
