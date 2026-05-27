#!/bin/bash

# 配置
DB_HOST="localhost"
DB_PORT="3306"
DB_USER="dsm_user"
DB_PASSWORD="asdf3asRDSfEre4DAS79"
DB_NAME="dsm"

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