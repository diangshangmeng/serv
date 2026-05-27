# 数据库导入指南

## 概述

本文档说明如何使用 `schema.sql` 文件初始化数据库。

## 前置条件

- MySQL 5.7+ 或 MariaDB 10.2+
- MySQL 客户端工具（mysql CLI）或图形化工具（如 MySQL Workbench、Navicat）

## 文件位置

```bash
database/schema.sql
```

## 导入步骤

### 方法一：使用 MySQL 命令行（推荐）

#### 1. 创建数据库

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE voucher_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 退出
EXIT;
```

#### 2. 导入 Schema

```bash
# 方式1：直接导入
mysql -u root -p voucher_platform < database/schema.sql

# 方式2：如果有错误，查看详细输出
mysql -u root -p voucher_platform -e "source database/schema.sql"

# 方式3：交互式导入
mysql -u root -p voucher_platform
mysql> source database/schema.sql
```

#### 3. 验证导入结果

```bash
# 登录数据库
mysql -u root -p voucher_platform

# 查看所有表
SHOW TABLES;

# 查看城市数据
SELECT * FROM cities;

# 查看管理员
SELECT * FROM admins;

# 退出
EXIT;
```

### 方法二：使用 Docker

```bash
# 如果使用 Docker 运行 MySQL
docker exec -it <container_name> mysql -u root -p voucher_platform < database/schema.sql
```

### 方法三：使用 GUI 工具

1. **MySQL Workbench**
   - 打开 MySQL Workbench
   - 连接到数据库服务器
   - 选择 "File" → "Run SQL Script"
   - 选择 `database/schema.sql` 文件
   - 执行导入

2. **Navicat**
   - 连接到数据库服务器
   - 创建数据库 `voucher_platform`
   - 右键点击数据库 → "Execute SQL File"
   - 选择 `database/schema.sql` 文件
   - 执行导入

3. **DBeaver**
   - 连接到数据库服务器
   - 右键点击数据库 → "SQL" → "Execute Script"
   - 选择 `database/schema.sql` 文件
   - 执行导入

## 导入内容

### 数据表列表

| 序号 | 表名 | 说明 |
|------|------|------|
| 1 | cities | 城市表 |
| 2 | users | 用户表 |
| 3 | admins | 管理员表 |
| 4 | product_images | 产品图片表 |
| 5 | payment_code_images | 支付码图片表 |
| 6 | products | 产品表 |
| 7 | coupons | 券码表 |
| 8 | listings | 挂牌表 |
| 9 | orders | 订单表 |

### 初始数据

#### 城市数据（云南省）
- 昆明市
- 大理市
- 曲靖市
- 玉溪市
- 保山市
- 昭通市
- 丽江市
- 普洱市
- 临沧市

#### 管理员账户
- 用户名：`admin`
- 密码：`admin123`
- 密码哈希：`$2a$10$N9qo8uLOickgx2ZMRZoMye.IjzqAKL9xL5jvMFVdNJHvGCgTq/VEq`

#### 测试用户账户
- 用户1：手机号 `13800138001`，密码 `123456`，未认证
- 用户2：手机号 `13900139001`，密码 `123456`，待审核
- 用户3：手机号 `13700137001`，密码 `123456`，已认证

## 常见问题

### Q1: 导入时报错 "Unknown database"
**原因**：数据库不存在
**解决**：先执行 `CREATE DATABASE voucher_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;`

### Q2: 导入时报错 "Table already exists"
**原因**：`schema.sql` 包含 DROP TABLE 语句，检查是否重复导入
**解决**：如果是重复导入，这是正常现象。如果不想删除现有数据，可以手动导入单个表或修改 SQL 文件移除 DROP TABLE 语句

### Q3: 字符集问题
**原因**：数据库字符集不是 utf8mb4
**解决**：确保创建数据库时指定正确的字符集
```sql
CREATE DATABASE voucher_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### Q4: 权限问题
**原因**：MySQL 用户没有创建表的权限
**解决**：使用有权限的用户登录，或联系 DBA 授权

## 一键导入脚本

### Linux/macOS

创建 `import_db.sh` 脚本：

```bash
#!/bin/bash

# 配置
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="root"
DB_PASSWORD="your_password"
DB_NAME="voucher_platform"

echo "开始导入数据库..."

# 创建数据库
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD -e \
"CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 导入数据
mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME < database/schema.sql

if [ $? -eq 0 ]; then
    echo "数据库导入成功！"
else
    echo "数据库导入失败！"
    exit 1
fi
```

使用方式：
```bash
chmod +x import_db.sh
./import_db.sh
```

### Windows

创建 `import_db.bat` 脚本：

```batch
@echo off
set DB_HOST=localhost
set DB_PORT=3306
set DB_USER=root
set DB_PASSWORD=your_password
set DB_NAME=voucher_platform

echo 开始导入数据库...

mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% -e "CREATE DATABASE IF NOT EXISTS %DB_NAME% CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

mysql -h %DB_HOST% -P %DB_PORT% -u %DB_USER% -p%DB_PASSWORD% %DB_NAME% < database\schema.sql

if %errorlevel% equ 0 (
    echo 数据库导入成功！
) else (
    echo 数据库导入失败！
    pause
    exit /b 1
)

pause
```

## 验证清单

导入完成后，请验证以下内容：

- [ ] 所有 9 个数据表已创建
- [ ] 城市表包含 9 条记录
- [ ] 管理员表包含 1 条记录
- [ ] 用户表包含 3 条测试用户记录
- [ ] 所有外键关系正确（如果有的话）

## 后续步骤

1. 修改配置文件 `config.yaml` 中的数据库连接信息
2. 启动应用程序
3. 使用管理员账号登录后台管理系统
4. 测试各项功能

## 技术支持

如有疑问，请联系开发团队。
