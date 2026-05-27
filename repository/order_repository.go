package repository

import (
	"time"
	"voucher-platform/config"
	"voucher-platform/model"
)

func CreateOrder(order *model.Order) error {
	return config.DB.Create(order).Error
}

func GetOrderByID(id uint) (*model.Order, error) {
	var order model.Order
	err := config.DB.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func GetOrderByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := config.DB.Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func GetOrdersByBuyerID(phone string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{}).Where("buyer_phone = ?", phone)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetOrdersBySellerPhone(phone string, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{}).Where("seller_phone = ?", phone)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func UpdateOrder(order *model.Order) error {
	return config.DB.Save(order).Error
}

func GetExpiredOrders() ([]model.Order, error) {
	var orders []model.Order
	err := config.DB.Where("status = ? AND expired_at < ?", 0, time.Now()).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func UpdateOrderStatus(orderID uint, status int) error {
	return config.DB.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func GetBuyerOrdersByPhone(phone string, status *int, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{}).Where("buyer_phone = ?", phone)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetSellerOrdersByPhone(phone string, status *int, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{}).Where("seller_phone = ?", phone)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetAllOrdersByPhone(phone string, status *int, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{}).Where("buyer_phone = ? OR seller_phone = ?", phone, phone)

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func AdminGetOrders(orderNo, phone string, status *int, startTime, endTime *time.Time, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{})

	if orderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+orderNo+"%")
	}

	if phone != "" {
		query = query.Where("buyer_phone LIKE ? OR seller_phone LIKE ?", "%"+phone+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}

	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func GetOrderStats() (map[string]int64, error) {
	stats := make(map[string]int64)
	var total int64
	
	err := config.DB.Model(&model.Order{}).Count(&total).Error
	if err != nil {
		return nil, err
	}
	stats["total"] = total

	var unpaid, paid, completed, cancelled, timeout int64
	err = config.DB.Model(&model.Order{}).Where("status = ?", 0).Count(&unpaid).Error
	if err != nil {
		return nil, err
	}
	stats["unpaid"] = unpaid

	err = config.DB.Model(&model.Order{}).Where("status = ?", 1).Count(&paid).Error
	if err != nil {
		return nil, err
	}
	stats["paid"] = paid

	err = config.DB.Model(&model.Order{}).Where("status = ?", 2).Count(&completed).Error
	if err != nil {
		return nil, err
	}
	stats["completed"] = completed

	err = config.DB.Model(&model.Order{}).Where("status = ?", 3).Count(&cancelled).Error
	if err != nil {
		return nil, err
	}
	stats["cancelled"] = cancelled

	err = config.DB.Model(&model.Order{}).Where("status = ?", 4).Count(&timeout).Error
	if err != nil {
		return nil, err
	}
	stats["timeout"] = timeout

	return stats, nil
}

func GetActiveOrderByProductID(productID uint64) (*model.Order, error) {
	var order model.Order
	err := config.DB.Where("product_id = ? AND status IN ?", productID, []int{0, 1}).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
