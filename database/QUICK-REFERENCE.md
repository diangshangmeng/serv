# 数据库 DDL 快速参考

## 一键导入命令

### 基本导入
```bash
mysql -u root -p voucher_platform < database/schema.sql
```

### 完整流程
```bash
# 1. 创建数据库
mysql -u root -p -e "CREATE DATABASE voucher_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 2. 导入数据
mysql -u root -p voucher_platform < database/schema.sql

# 3. 验证
mysql -u root -p voucher_platform -e "SHOW TABLES;"
```

## 测试账号

| 类型 | 账号 | 密码 |
|------|------|------|
| 管理员 | admin | admin123 |
| 用户 | 13800138001 | 123456 |
| 用户 | 13900139001 | 123456 |
| 用户 | 13700137001 | 123456 |

## 数据表

1. `cities` - 城市表
2. `users` - 用户表
3. `admins` - 管理员表
4. `product_images` - 产品图片表
5. `payment_code_images` - 支付码图片表
6. `products` - 产品表
7. `coupons` - 券码表
8. `listings` - 挂牌表
9. `orders` - 订单表

## 初始数据

- 9 个城市（云南省）
- 1 个管理员
- 3 个测试用户

详细说明请查看：[database/README.md](file:///Users/richard/workspace/1.MerCoupontX/serv/database/README.md)
