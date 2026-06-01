package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/util"
)

const (
	productCachePrefix       = "product:"
	productListCachePrefix   = "list:"
	productDetailCachePrefix = "detail:"
	myProductsCachePrefix    = "my_products:"
	productListCacheTTL      = 5 * time.Minute
	productDetailCacheTTL    = 10 * time.Minute
	myProductsCacheTTL       = 5 * time.Minute
)

// 定义一个包装商品列表和总数的结构体
type ProductListCache struct {
	List  []model.Product `json:"list"`
	Total int64           `json:"total"`
}

type MyProductsListCache struct {
	List  []*model.Product `json:"list"`
	Total int64            `json:"total"`
}

// GetProductListFromCache 从Redis缓存获取商品列表
func GetProductListFromCache(page int, pageSize int, cityID *uint) (*ProductListCache, error) {
	ctx := context.Background()
	cacheKey := generateProductListCacheKey(page, pageSize, cityID)

	data, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var cache ProductListCache
	err = json.Unmarshal([]byte(data), &cache)
	if err != nil {
		util.GetLogger().Error("product_cache_unmarshal_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
		return nil, err
	}

	return &cache, nil
}

// SetProductListToCache 将商品列表写入Redis缓存
func SetProductListToCache(list []model.Product, total int64, page int, pageSize int, cityID *uint) error {
	ctx := context.Background()
	cacheKey := generateProductListCacheKey(page, pageSize, cityID)

	cache := ProductListCache{
		List:  list,
		Total: total,
	}

	data, err := json.Marshal(cache)
	if err != nil {
		util.GetLogger().Error("product_cache_marshal_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
		return err
	}

	err = config.RedisClient.Set(ctx, cacheKey, data, productListCacheTTL).Err()
	if err != nil {
		util.GetLogger().Error("product_cache_set_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
	}

	return err
}

// GetProductDetailFromCache 从Redis缓存获取商品详情
func GetProductDetailFromCache(productID uint64) (*model.Product, error) {
	ctx := context.Background()
	cacheKey := generateProductDetailCacheKey(productID)

	data, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var product model.Product
	err = json.Unmarshal([]byte(data), &product)
	if err != nil {
		util.GetLogger().Error("product_detail_cache_unmarshal_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
		return nil, err
	}

	return &product, nil
}

// SetProductDetailToCache 将商品详情写入Redis缓存
func SetProductDetailToCache(product *model.Product) error {
	ctx := context.Background()
	cacheKey := generateProductDetailCacheKey(uint64(product.ID))

	data, err := json.Marshal(product)
	if err != nil {
		util.GetLogger().Error("product_detail_cache_marshal_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
		return err
	}

	err = config.RedisClient.Set(ctx, cacheKey, data, productDetailCacheTTL).Err()
	if err != nil {
		util.GetLogger().Error("product_detail_cache_set_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
	}

	return err
}

// ClearProductDetailCache 清除商品详情缓存
func ClearProductDetailCache(productID uint64) {
	ctx := context.Background()
	cacheKey := generateProductDetailCacheKey(productID)

	err := config.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		util.GetLogger().Warn("product_detail_cache_clear_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
	}
}

// ClearProductListCache 清除商品列表缓存（使用模糊匹配清除所有列表缓存）
func ClearProductListCache() {
	ctx := context.Background()
	pattern := productCachePrefix + productListCachePrefix + "*"

	keys, err := scanRedisKeys(ctx, pattern)
	if err != nil {
		util.GetLogger().Warn("product_list_cache_scan_failed", util.StringField("error", err.Error()))
		return
	}

	if len(keys) > 0 {
		err = config.RedisClient.Del(ctx, keys...).Err()
		if err != nil {
			util.GetLogger().Warn("product_list_cache_clear_failed", util.StringField("error", err.Error()))
		}
	}
}

// ClearAllProductCache 清除所有商品相关的缓存
func ClearAllProductCache() {
	ClearProductListCache()
}

// generateProductListCacheKey 生成商品列表缓存键
func generateProductListCacheKey(page int, pageSize int, cityID *uint) string {
	if cityID != nil {
		return fmt.Sprintf("%s%s%d:%d:%d", productCachePrefix, productListCachePrefix, *cityID, page, pageSize)
	}
	return fmt.Sprintf("%s%sall:%d:%d", productCachePrefix, productListCachePrefix, page, pageSize)
}

// generateProductDetailCacheKey 生成商品详情缓存键
func generateProductDetailCacheKey(productID uint64) string {
	return fmt.Sprintf("%s%s%d", productCachePrefix, productDetailCachePrefix, productID)
}

func GetMyProductsFromCache(phone string, page int, pageSize int, cityID *uint) (*MyProductsListCache, error) {
	ctx := context.Background()
	cacheKey := generateMyProductsCacheKey(phone, page, pageSize, cityID)

	data, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var cache MyProductsListCache
	err = json.Unmarshal([]byte(data), &cache)
	if err != nil {
		util.GetLogger().Error("my_products_cache_unmarshal_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
		return nil, err
	}

	return &cache, nil
}

func SetMyProductsToCache(list []*model.Product, total int64, phone string, page int, pageSize int, cityID *uint) error {
	ctx := context.Background()
	cacheKey := generateMyProductsCacheKey(phone, page, pageSize, cityID)

	cache := MyProductsListCache{
		List:  list,
		Total: total,
	}

	data, err := json.Marshal(cache)
	if err != nil {
		util.GetLogger().Error("my_products_cache_marshal_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
		return err
	}

	err = config.RedisClient.Set(ctx, cacheKey, data, myProductsCacheTTL).Err()
	if err != nil {
		util.GetLogger().Error("my_products_cache_set_failed", util.StringField("key", cacheKey), util.StringField("error", err.Error()))
	}

	return err
}

func ClearMyProductsCache(phone string) {
	ctx := context.Background()
	pattern := productCachePrefix + myProductsCachePrefix + phone + ":*"

	keys, err := scanRedisKeys(ctx, pattern)
	if err != nil {
		util.GetLogger().Warn("my_products_cache_scan_failed", util.StringField("error", err.Error()))
		return
	}

	if len(keys) > 0 {
		err = config.RedisClient.Del(ctx, keys...).Err()
		if err != nil {
			util.GetLogger().Warn("my_products_cache_clear_failed", util.StringField("error", err.Error()))
		}
	}
}

func generateMyProductsCacheKey(phone string, page int, pageSize int, cityID *uint) string {
	if cityID != nil {
		return fmt.Sprintf("%s%s%s:%d:%d:%d", productCachePrefix, myProductsCachePrefix, phone, *cityID, page, pageSize)
	}
	return fmt.Sprintf("%s%s%s:all:%d:%d", productCachePrefix, myProductsCachePrefix, phone, page, pageSize)
}
