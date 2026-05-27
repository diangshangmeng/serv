# 数据库初始化说明

## 快速开始

### 方法一：使用 SQL 文件（推荐）

1. **登录 MySQL**:
   ```bash
   mysql -u root -p
   ```

2. **执行 SQL 文件**:
   ```bash
   mysql -u root -p < migrations/init_database.sql
   ```

   或者在 MySQL 命令行中：
   ```sql
   SOURCE /path/to/migrations/init_database.sql;
   ```

### 方法二：使用 GORM AutoMigrate（开发环境）

应用程序启动时会自动执行 `model.AutoMigrate()`，会自动创建所有表。

```go
// 在 main.go 中会自动调用
err = model.AutoMigrate()
```

## 数据库表结构

### 1. users（用户表）
- `phone`: 手机号（唯一）
- `status`: 状态（0封禁，1正常）
- `auth_status`: 认证状态（0未认证，1待审核，2已认证，3已驳回）

### 2. cities（城市表）
- 云南省的9个城市

### 3. coupons（代金券表）
- 用户持有的代金券

### 4. listings（上架表）
- 代金券上架记录

### 5. orders（订单表）- 核心表
- **订单号**: `order_no`（雪花ID，唯一）
- **商品快照**: 标题、描述、图片等
- **买卖双方**: seller_id, buyer_id, 电话
- **状态**: 0未付款，1已付款待确认，2交易完成，3交易取消，4超时关闭
- **金额**: `price`（单位：分）
- **过期时间**: `expired_at`（15分钟后自动关闭）

### 6. admins（管理员表）
- 默认账号: admin / admin123

### 7. product_images（商品图片表）
- 商品图片管理

### 8. payment_code_images（付款码图片表）
- 付款码图片管理
- `company_name`: 公司名称

### 9. products（商品表）
- 商品管理
- **状态**: 0未上架，1上架，2锁定，4等待确认

## 测试账号

### 管理员账号
- **用户名**: admin
- **密码**: admin123
- **权限**: 后台管理权限

### 测试用户
| 手机号 | 密码 | 认证状态 | 说明 |
|--------|------|----------|------|
| 13800138001 | 123456 | 未认证 | 未提交认证 |
| 13900139001 | 123456 | 待审核 | 已提交认证，待审核 |
| 13700137001 | 123456 | 已认证 | 认证通过，可正常使用 |

## 验证安装

执行以下 SQL 查看数据：

```sql
-- 查看所有表
SHOW TABLES;

-- 查看用户
SELECT id, phone, status, auth_status FROM users;

-- 查看商品
SELECT id, title, status FROM products;

-- 查看订单
SELECT id, order_no, status, price FROM orders;
```

## 常见问题

### Q: 表已存在但字段不完整？
A: 手动执行 ALTER TABLE 语句添加缺失的字段。

### Q: 初始数据没有插入？
A: 使用 INSERT IGNORE 或先删除再插入。

### Q: 订单表有 shipping_address_id 错误？
A: 这是旧字段，已通过迁移脚本删除。如果数据库中仍有此字段，手动删除：
```sql
ALTER TABLE orders DROP COLUMN shipping_address_id;
```

## 重要提示

1. **代金券交易不需要地址**: orders 表不包含地址相关字段
2. **雪花ID生成订单号**: 确保 Snowflake 服务正常
3. **时区设置**: MySQL 使用 datetime(3) 毫秒级时间戳

## 技术细节

- **字符集**: utf8mb4（支持emoji和特殊字符）
- **存储引擎**: InnoDB（支持事务）
- **软删除**: 所有表都包含 `deleted_at` 字段
- **时间戳**: 使用 datetime(3) 毫秒级精度
