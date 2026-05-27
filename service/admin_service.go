package service

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"voucher-platform/config"
	"voucher-platform/model"
	"voucher-platform/repository"
	"voucher-platform/util"
)

func AdminLogin(username, password string) (string, error) {
	admin, err := repository.GetAdminByUsername(username)
	if err != nil {
		return "", util.NewBizError(util.ErrCodeAdminAuthFailed, "用户名或密码错误")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", util.NewBizError(util.ErrCodeAdminAuthFailed, "用户名或密码错误")
	}

	return util.GenerateAdminToken(uint64(admin.ID), admin.Username)
}

func GetPendingUsers() ([]model.User, error) {
	users, err := repository.GetPendingUsers()
	if err != nil {
		return nil, err
	}
	ConvertUsersImageURLs(users)
	return users, nil
}

func GetAllUsers() ([]model.User, error) {
	users, err := repository.GetAllUsers()
	if err != nil {
		return nil, err
	}
	ConvertUsersImageURLs(users)
	return users, nil
}

func GetAllOrders(page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := config.DB.Model(&model.Order{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	ConvertOrdersImageURLs(orders)

	return orders, total, nil
}

func GetUserDetail(userID uint64) (*model.User, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}
	ConvertUserImageURLs(user)
	return user, nil
}

func SearchUsers(phone string, status *int, page, pageSize int) ([]model.User, int64, error) {
	users, total, err := repository.GetUsersWithFilters(phone, status, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	ConvertUsersImageURLs(users)
	return users, total, nil
}

func BanUser(userID uint64) error {
	log.Printf("[BanUser] 开始封禁用户，user_id=%d", userID)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Printf("[BanUser] 用户不存在，user_id=%d, err=%v", userID, err)
		return util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}

	user.Status = 0
	log.Printf("[BanUser] 用户已封禁，user_id=%d", userID)
	return repository.UpdateUser(user)
}

func UnbanUser(userID uint64) error {
	log.Printf("[UnbanUser] 开始解封用户，user_id=%d", userID)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Printf("[UnbanUser] 用户不存在，user_id=%d, err=%v", userID, err)
		return util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}

	user.Status = 1
	log.Printf("[UnbanUser] 用户已解封，user_id=%d", userID)
	return repository.UpdateUser(user)
}

func AuditUser(userID uint64, pass bool, remark string) error {
	log.Printf("[AuditUser] 开始审核用户，user_id=%d, pass=%t, remark=%s", userID, pass, remark)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Printf("[AuditUser] 用户不存在，user_id=%d, err=%v", userID, err)
		return util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}

	if pass {
		user.AuthStatus = model.AuthStatusApproved
	} else {
		user.AuthStatus = model.AuthStatusRejected
	}
	user.AuditRemark = remark
	log.Printf("[AuditUser] 用户审核完成，user_id=%d, status=%d", userID, user.AuthStatus)
	return repository.UpdateUser(user)
}

func SetUserStatus(userID uint64, status int) error {
	log.Printf("[SetUserStatus] 开始设置用户状态，user_id=%d, status=%d", userID, status)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Printf("[SetUserStatus] 用户不存在，user_id=%d, err=%v", userID, err)
		return util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}

	user.Status = status
	log.Printf("[SetUserStatus] 用户状态已更新，user_id=%d, status=%d", userID, user.Status)
	return repository.UpdateUser(user)
}

func ResetUserPassword(userID uint64) (string, error) {
	log.Printf("[ResetUserPassword] 开始重置用户密码，user_id=%d", userID)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Printf("[ResetUserPassword] 用户不存在，user_id=%d, err=%v", userID, err)
		return "", util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}

	newPassword := "123456"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ResetUserPassword] 密码加密失败，user_id=%d, err=%v", userID, err)
		return "", err
	}

	user.Password = string(hashedPassword)
	err = repository.UpdateUser(user)
	if err != nil {
		log.Printf("[ResetUserPassword] 更新用户失败，user_id=%d, err=%v", userID, err)
		return "", err
	}

	log.Printf("[ResetUserPassword] 用户密码重置成功，user_id=%d", userID)
	return newPassword, nil
}

func UpdateUser(userID uint64, cityID *uint, status *int, authStatus *int, auditRemark string, idCardFront string, idCardBack string, businessLicense string, paymentCode string) (*model.User, error) {
	log.Printf("[UpdateUser] 开始更新用户信息，user_id=%d", userID)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		log.Printf("[UpdateUser] 用户不存在，user_id=%d, err=%v", userID, err)
		return nil, util.NewBizError(util.ErrCodeUserNotFound, "用户不存在")
	}

	if cityID != nil {
		user.CityID = *cityID
	}
	if status != nil {
		user.Status = *status
	}
	if authStatus != nil {
		user.AuthStatus = *authStatus
	}
	if auditRemark != "" {
		user.AuditRemark = auditRemark
	}
	if idCardFront != "" {
		user.IDCardFront = idCardFront
	}
	if idCardBack != "" {
		user.IDCardBack = idCardBack
	}
	if businessLicense != "" {
		user.BusinessLicense = businessLicense
	}
	if paymentCode != "" {
		user.PaymentCode = paymentCode
	}

	err = repository.UpdateUser(user)
	if err != nil {
		log.Printf("[UpdateUser] 更新用户失败，user_id=%d, err=%v", userID, err)
		return nil, err
	}

	log.Printf("[UpdateUser] 用户信息更新成功，user_id=%d", userID)
	ConvertUserImageURLs(user)
	return user, nil
}

func AdminSellerCancelOrder(orderNo string) error {
	log.Printf("[AdminSellerCancelOrder] 开始管理员取消订单，order_no=%s", orderNo)
	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		log.Printf("[AdminSellerCancelOrder] 订单不存在，order_no=%s, err=%v", orderNo, err)
		return util.NewBizError(util.ErrCodeOrderNotFound, "订单不存在")
	}

	if order.Status != ORDER_STATUS_PAID {
		log.Printf("[AdminSellerCancelOrder] 订单状态不正确，order_no=%s, status=%d", orderNo, order.Status)
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确，只能取消已付款待确认的订单")
	}

	tx := config.DB.Begin()

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

	now := time.Now()
	order.Status = ORDER_STATUS_UNPAID
	order.ExpiredAt = now.Add(15 * time.Minute)
	order.PaymentVoucher = ""
	err = tx.Save(order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	log.Printf("[AdminSellerCancelOrder] 管理员取消订单成功，order_no=%s", orderNo)
	return nil
}

func AdminSellerConfirmOrder(orderNo string) error {
	log.Printf("[AdminSellerConfirmOrder] 开始管理员确认订单，order_no=%s", orderNo)
	order, err := repository.GetOrderByOrderNo(orderNo)
	if err != nil {
		log.Printf("[AdminSellerConfirmOrder] 订单不存在，order_no=%s, err=%v", orderNo, err)
		return util.NewBizError(util.ErrCodeOrderNotFound, "订单不存在")
	}

	if order.Status != ORDER_STATUS_PAID {
		log.Printf("[AdminSellerConfirmOrder] 订单状态不正确，order_no=%s, status=%d", orderNo, order.Status)
		return util.NewBizError(util.ErrCodeOrderStatusError, "订单状态不正确，只能确认已付款待确认的订单")
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
				log.Printf("[AdminSellerConfirmOrder] 更新 coupon 失败，继续执行，order_no=%s, err=%v", orderNo, err)
			}
		} else {
			log.Printf("[AdminSellerConfirmOrder] 查询 coupon 失败，继续执行，order_no=%s, err=%v", orderNo, err)
		}
	} else {
		log.Printf("[AdminSellerConfirmOrder] 查询 listing 失败，继续执行，order_no=%s, err=%v", orderNo, err)
	}

	product, err := repository.GetProductByID(uint64(order.ProductID))
	if err != nil {
		tx.Rollback()
		return err
	}
	product.Status = 0
	product.OwnerID = order.BuyerID
	product.OwnerPhone = order.BuyerPhone
	
	buyer, err := repository.GetUserByID(uint64(order.BuyerID))
	if err != nil {
		tx.Rollback()
		return err
	}
	product.PaymentImageURL = buyer.PaymentCode
	
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

	tx.Commit()
	log.Printf("[AdminSellerConfirmOrder] 管理员确认订单成功，order_no=%s", orderNo)
	return nil
}
