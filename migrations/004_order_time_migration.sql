-- 迁移脚本：订单时间字段迁移
-- 日期: 2026-05-18
-- 描述: 将order_time字段从products表迁移到orders表

-- 1. 向orders表添加order_time列
ALTER TABLE orders ADD COLUMN order_time DATETIME(3) DEFAULT NULL COMMENT '下单时间' AFTER expired_at;

-- 2. 从products表删除order_time列（如果存在）
ALTER TABLE products DROP COLUMN IF EXISTS order_time;

-- 3. 验证迁移结果
-- SELECT order_time, expired_at FROM orders LIMIT 10;
-- DESCRIBE orders;
-- DESCRIBE products;
