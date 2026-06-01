-- ============================================
-- 代金券交易平台 - 数据库初始化脚本
-- 创建日期: 2026-05-17
-- 数据库版本: MySQL 8.0+
-- ============================================

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS voucher_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE voucher_platform;

-- ============================================
-- 1. 用户表 (users)
-- ============================================
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `phone` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户手机号',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码（加密）',
  `city_id` bigint unsigned DEFAULT NULL COMMENT '城市ID',
  `status` int DEFAULT '1' COMMENT '状态：0封禁，1正常',
  `auth_status` int DEFAULT '0' COMMENT '认证状态：0未认证，1待审核，2已认证，3已驳回',
  `audit_remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '审核备注',
  `id_card_front` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '身份证正面',
  `id_card_back` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '身份证背面',
  `business_license` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '营业执照',
  `payment_code` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '付款码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ============================================
-- 2. 城市表 (cities)
-- ============================================
CREATE TABLE IF NOT EXISTS `cities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '城市名称',
  `province_id` bigint unsigned DEFAULT '1' COMMENT '省份ID',
  PRIMARY KEY (`id`),
  KEY `idx_cities_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='城市表';

-- ============================================
-- 3. 代金券表 (coupons)
-- ============================================
CREATE TABLE IF NOT EXISTS `coupons` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `code` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '代金券码',
  `original_price` double NOT NULL COMMENT '原价',
  `status` int DEFAULT '0' COMMENT '状态：0未使用，1已使用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_coupons_code` (`code`),
  KEY `idx_coupons_user_id` (`user_id`),
  KEY `idx_coupons_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='代金券表';

-- ============================================
-- 4. 上架表 (listings)
-- ============================================
CREATE TABLE IF NOT EXISTS `listings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `coupon_id` bigint unsigned NOT NULL COMMENT '代金券ID',
  `product_id` bigint unsigned NOT NULL COMMENT '商品ID',
  `selling_price` double NOT NULL COMMENT '出售价格',
  `status` int DEFAULT '0' COMMENT '状态：0待售，1已售',
  PRIMARY KEY (`id`),
  KEY `idx_listings_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='上架表';

-- ============================================
-- 5. 订单表 (orders) - 核心交易表
-- ============================================
CREATE TABLE IF NOT EXISTS `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_no` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单号（雪花ID）',
  `product_id` bigint unsigned DEFAULT NULL COMMENT '商品ID',
  `product_title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品标题（快照）',
  `product_description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品描述（快照）',
  `product_city` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品城市（快照）',
  `seller_id` bigint unsigned DEFAULT NULL COMMENT '卖家ID',
  `seller_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '卖家手机号',
  `product_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品图片URL（快照）',
  `payment_code_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '付款码图片URL（快照）',
  `buyer_id` bigint unsigned DEFAULT NULL COMMENT '买家ID',
  `buyer_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '买家手机号',
  `payment_voucher` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '付款凭证URL',
  `status` int DEFAULT '0' COMMENT '订单状态：0未付款，1已付款待确认，2交易完成，3交易取消，4超时关闭',
  `price` bigint DEFAULT NULL COMMENT '订单金额（单位：分）',
  `expired_at` datetime(3) DEFAULT NULL COMMENT '订单过期时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_no` (`order_no`),
  KEY `idx_orders_product_id` (`product_id`),
  KEY `idx_orders_buyer_id` (`buyer_id`),
  KEY `idx_orders_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- ============================================
-- 6. 管理员表 (admins)
-- ============================================
CREATE TABLE IF NOT EXISTS `admins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '管理员用户名',
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码（加密）',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  KEY `idx_admins_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- ============================================
-- 7. 商品图片表 (product_images)
-- ============================================
CREATE TABLE IF NOT EXISTS `product_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图片名称',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `tags` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标签',
  `is_used` tinyint(1) DEFAULT '0' COMMENT '是否已使用',
  `admin_id` bigint unsigned DEFAULT NULL COMMENT '管理员ID',
  `size` bigint DEFAULT NULL COMMENT '文件大小（字节）',
  PRIMARY KEY (`id`),
  KEY `idx_product_images_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品图片表';

-- ============================================
-- 8. 付款码图片表 (payment_code_images)
-- ============================================
CREATE TABLE IF NOT EXISTS `payment_code_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '存储路径',
  `company_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '公司名称',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `admin_id` bigint unsigned DEFAULT NULL COMMENT '管理员ID',
  `size` bigint DEFAULT NULL COMMENT '文件大小（字节）',
  PRIMARY KEY (`id`),
  KEY `idx_payment_code_images_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='付款码图片表';

-- ============================================
-- 9. 商品表 (products)
-- ============================================
CREATE TABLE IF NOT EXISTS `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品标题',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品描述',
  `price` bigint DEFAULT NULL COMMENT '商品价格（单位：分）',
  `city` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '商品城市',
  `display_image_id` bigint unsigned DEFAULT NULL COMMENT '展示图片ID',
  `payment_image_id` bigint unsigned DEFAULT NULL COMMENT '付款码图片ID',
  `owner_id` bigint unsigned DEFAULT NULL COMMENT '所有者ID',
  `owner_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所有者手机号',
  `display_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '展示图片URL',
  `payment_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '付款码图片URL',
  `status` int DEFAULT '0' COMMENT '商品状态：0未上架，1上架，2锁定，4等待确认',
  `order_time` datetime(3) DEFAULT NULL COMMENT '下单时间',
  `pay_time` datetime(3) DEFAULT NULL COMMENT '支付时间',
  `lock_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '锁定原因',
  `city_id` bigint unsigned DEFAULT NULL COMMENT '城市ID',
  `version` int DEFAULT '0' COMMENT '乐观锁版本号',
  `last_transaction_price` bigint DEFAULT '0' COMMENT '上次交易价格（分），用于防篡改',
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`),
  KEY `idx_products_owner_id` (`owner_id`),
  KEY `idx_products_status` (`status`),
  KEY `idx_products_city_id` (`city_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品表';

-- ============================================
-- 初始化数据
-- ============================================

-- 插入城市数据
INSERT IGNORE INTO `cities` (`id`, `name`, `province_id`) VALUES
(1, '昆明市', 1),
(2, '大理市', 1),
(3, '曲靖市', 1),
(4, '玉溪市', 1),
(5, '保山市', 1),
(6, '昭通市', 1),
(7, '丽江市', 1),
(8, '普洱市', 1),
(9, '临沧市', 1);

-- 插入管理员账号
-- 用户名: admin
-- 密码: admin123 (BCRYPT加密)
INSERT IGNORE INTO `admins` (`id`, `username`, `password`) VALUES
(1, 'admin', '$2a$10$kj1dClEJF9gbL6axUER.PeYNajymob4cI9r/I9YxPiiIYsFjGyp4q');

-- ============================================
-- 完成提示
-- ============================================
SELECT '========================================' AS '';
SELECT '数据库初始化完成！' AS '恭喜';
SELECT '========================================' AS '';

-- 显示所有表
SHOW TABLES;

-- 显示用户数据
SELECT id, phone, city_id, status, auth_status FROM users;

-- 显示管理员数据
SELECT id, username FROM admins;

-- 显示城市数据
SELECT id, name, province_id FROM cities;
