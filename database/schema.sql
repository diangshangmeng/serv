-- =============================================
-- Voucher Platform Database Schema
-- Generated: 2026-05-27
-- Database: MySQL 5.7+
-- =============================================

-- Drop existing tables (in reverse order of dependencies)
DROP TABLE IF EXISTS `orders`;
DROP TABLE IF EXISTS `listings`;
DROP TABLE IF EXISTS `coupons`;
DROP TABLE IF EXISTS `products`;
DROP TABLE IF EXISTS `payment_code_images`;
DROP TABLE IF EXISTS `product_images`;
DROP TABLE IF EXISTS `admins`;
DROP TABLE IF EXISTS `users`;
DROP TABLE IF EXISTS `cities`;

-- =============================================
-- Table 1: cities (城市表)
-- =============================================
CREATE TABLE `cities` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `province_id` bigint unsigned DEFAULT '1',
  PRIMARY KEY (`id`),
  KEY `idx_cities_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- Table 2: users (用户表)
-- =============================================
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `phone` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `city_id` bigint unsigned DEFAULT NULL COMMENT '用户所在城市ID',
  `status` int DEFAULT '1' COMMENT '用户状态：1-正常',
  `auth_status` int DEFAULT '0' COMMENT '认证状态：0-未认证，1-待审核，2-已认证，3-认证驳回',
  `audit_remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '审核备注',
  `id_card_front` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '身份证正面',
  `id_card_back` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '身份证背面',
  `business_license` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '营业执照',
  `payment_code` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '收款码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_phone` (`phone`),
  KEY `idx_users_deleted_at` (`deleted_at`),
  KEY `idx_users_city_id` (`city_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- =============================================
-- Table 3: admins (管理员表)
-- =============================================
CREATE TABLE `admins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_admins_username` (`username`),
  KEY `idx_admins_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员表';

-- =============================================
-- Table 4: product_images (产品图片表)
-- =============================================
CREATE TABLE `product_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图片路径',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图片名称',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `tags` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '标签',
  `is_used` tinyint(1) DEFAULT '0' COMMENT '是否已使用',
  `admin_id` bigint unsigned DEFAULT NULL COMMENT '上传管理员ID',
  `size` bigint DEFAULT NULL COMMENT '文件大小(字节)',
  PRIMARY KEY (`id`),
  KEY `idx_product_images_deleted_at` (`deleted_at`),
  KEY `idx_product_images_admin_id` (`admin_id`),
  KEY `idx_product_images_is_used` (`is_used`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品图片表';

-- =============================================
-- Table 5: payment_code_images (支付码图片表)
-- =============================================
CREATE TABLE `payment_code_images` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `path` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '图片路径',
  `company_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '公司名称',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '文件名',
  `admin_id` bigint unsigned DEFAULT NULL COMMENT '上传管理员ID',
  `size` bigint DEFAULT NULL COMMENT '文件大小(字节)',
  PRIMARY KEY (`id`),
  KEY `idx_payment_code_images_deleted_at` (`deleted_at`),
  KEY `idx_payment_code_images_admin_id` (`admin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='支付码图片表';

-- =============================================
-- Table 6: products (产品表)
-- =============================================
CREATE TABLE `products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '产品标题',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  `price` bigint DEFAULT NULL COMMENT '产品价格(分)',
  `display_image_id` bigint unsigned DEFAULT NULL COMMENT '展示图片ID',
  `payment_image_id` bigint unsigned DEFAULT NULL COMMENT '支付凭证图片ID',
  `owner_id` bigint unsigned DEFAULT NULL COMMENT '拥有者用户ID',
  `owner_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '拥有者手机号',
  `display_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '展示图片URL',
  `payment_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '支付凭证图片URL',
  `status` int DEFAULT '0' COMMENT '产品状态：0-未上架，1-已上架，2-已锁定',
  `pay_time` datetime(3) DEFAULT NULL COMMENT '支付时间',
  `lock_reason` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '锁定原因',
  `city_id` bigint unsigned DEFAULT NULL COMMENT '城市ID',
  PRIMARY KEY (`id`),
  KEY `idx_products_deleted_at` (`deleted_at`),
  KEY `idx_products_city_id` (`city_id`),
  KEY `idx_products_owner_id` (`owner_id`),
  KEY `idx_products_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='产品表';

-- =============================================
-- Table 7: coupons (券码表)
-- =============================================
CREATE TABLE `coupons` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `code` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '券码',
  `original_price` double NOT NULL COMMENT '原价',
  `status` int DEFAULT '0' COMMENT '状态：0-未使用，1-已使用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_coupons_code` (`code`),
  KEY `idx_coupons_deleted_at` (`deleted_at`),
  KEY `idx_coupons_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='券码表';

-- =============================================
-- Table 8: listings (挂牌表)
-- =============================================
CREATE TABLE `listings` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `coupon_id` bigint unsigned NOT NULL COMMENT '券码ID',
  `product_id` bigint unsigned NOT NULL COMMENT '产品ID',
  `selling_price` double NOT NULL COMMENT '挂牌价格',
  `status` int DEFAULT '0' COMMENT '状态：0-挂牌中，1-已售出，2-已撤销',
  PRIMARY KEY (`id`),
  KEY `idx_listings_deleted_at` (`deleted_at`),
  KEY `idx_listings_coupon_id` (`coupon_id`),
  KEY `idx_listings_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='挂牌表';

-- =============================================
-- Table 9: orders (订单表)
-- =============================================
CREATE TABLE `orders` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `order_no` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单号',
  `product_id` bigint unsigned DEFAULT NULL COMMENT '产品ID',
  `product_title` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '产品标题',
  `product_description` text COLLATE utf8mb4_unicode_ci COMMENT '产品描述',
  `seller_id` bigint unsigned DEFAULT NULL COMMENT '卖家ID',
  `seller_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '卖家手机号',
  `product_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '产品图片URL',
  `payment_code_image_url` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '支付凭证图片URL',
  `buyer_id` bigint unsigned DEFAULT NULL COMMENT '买家ID',
  `buyer_phone` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '买家手机号',
  `payment_voucher` varchar(2048) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '支付凭证',
  `status` int DEFAULT '0' COMMENT '订单状态：0-待支付，1-已支付，2-已完成，3-已取消',
  `price` bigint DEFAULT NULL COMMENT '订单价格(分)',
  `expired_at` datetime DEFAULT NULL COMMENT '过期时间',
  `order_time` datetime(3) DEFAULT NULL COMMENT '下单时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_orders_order_no` (`order_no`),
  KEY `idx_orders_deleted_at` (`deleted_at`),
  KEY `idx_orders_product_id` (`product_id`),
  KEY `idx_orders_buyer_id` (`buyer_id`),
  KEY `idx_orders_seller_id` (`seller_id`),
  KEY `idx_orders_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='订单表';

-- =============================================
-- Insert Initial Data
-- =============================================

-- Insert cities (云南省城市)
INSERT INTO `cities` (`name`, `province_id`) VALUES
('昆明市', 1),
('大理市', 1),
('曲靖市', 1),
('玉溪市', 1),
('保山市', 1),
('昭通市', 1),
('丽江市', 1),
('普洱市', 1),
('临沧市', 1);

-- Insert admin user (password: admin123)
-- Password hash: $2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq
INSERT INTO `admins` (`username`, `password`) VALUES
('admin', '$2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq');

-- Insert test users (password: 123456 for all)
-- Password hash: $2a$10$OwnHkiQG1UPtf5AJO5qDeuV98lEUxfylLYAYoL4HBrPJJPEJvcxi
INSERT INTO `users` (`phone`, `password`, `city_id`, `status`, `auth_status`) VALUES
('13800138001', '$2a$10$OwnHkiQG1UPtf5AJO5qDeuV98lEUxfylLYAYoL4HBrPJJPEJvcxi', 1, 1, 0),
('13900139001', '$2a$10$OwnHkiQG1UPtf5AJO5qDeuV98lEUxfylLYAYoL4HBrPJJPEJvcxi', 2, 1, 1),
('13700137001', '$2a$10$OwnHkiQG1UPtf5AJO5qDeuV98lEUxfylLYAYoL4HBrPJJPEJvcxi', 3, 1, 2);

-- =============================================
-- Summary
-- =============================================
-- Total tables: 9
-- Initial data: 3 cities, 1 admin, 3 test users
--
-- Test credentials:
-- Admin: username=admin, password=admin123
-- Users: phone=13800138001/13900139001/13700137001, password=123456
