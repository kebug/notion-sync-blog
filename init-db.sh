#!/bin/bash

# Notion Sync Blog 数据库初始化脚本

echo "正在初始化 Notion Sync Blog 数据库..."

# 检查 MySQL 是否运行
if ! command -v mysql &> /dev/null; then
    echo "错误: 未找到 MySQL 客户端，请先安装 MySQL"
    exit 1
fi

# 读取配置
CONFIG_FILE="${1:-config.yaml}"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "错误: 配置文件 $CONFIG_FILE 不存在"
    exit 1
fi

# 从配置文件中提取数据库信息（使用 grep 和 awk 简化处理）
DB_HOST=$(grep -A1 "^database:" "$CONFIG_FILE" | grep "host:" | awk '{print $2}')
DB_PORT=$(grep -A1 "^database:" "$CONFIG_FILE" | grep "port:" | awk '{print $2}')
DB_USER=$(grep -A1 "^database:" "$CONFIG_FILE" | grep "user:" | awk '{print $2}')
DB_PASSWORD=$(grep -A1 "^database:" "$CONFIG_FILE" | grep "password:" | awk '{print $2}')
DB_NAME=$(grep -A1 "^database:" "$CONFIG_FILE" | grep "dbname:" | awk '{print $2}')

if [ -z "$DB_HOST" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    echo "错误: 无法从配置文件中解析数据库信息"
    exit 1
fi

# 执行数据库初始化
echo "连接到数据库: $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASSWORD" < backend/migrations/001_init.sql

if [ $? -eq 0 ]; then
    echo "数据库初始化成功!"
    echo "数据库: $DB_NAME 已创建"
    echo "表结构已导入"
else
    echo "错误: 数据库初始化失败"
    exit 1
fi
