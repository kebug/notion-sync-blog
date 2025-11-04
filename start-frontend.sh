#!/bin/bash

echo "🚀 启动 Notion Sync Blog 前端"

# 检查是否在 frontend 目录
if [ ! -f "frontend/package.json" ]; then
  echo "❌ 错误：未找到 frontend 目录"
  exit 1
fi

cd frontend

# 安装依赖（如果 node_modules 不存在）
if [ ! -d "node_modules" ]; then
  echo "📦 安装依赖..."
  npm install
fi

# 启动开发服务器
echo "✨ 启动开发服务器..."
echo "前端地址: http://localhost:3000"
echo "API 地址: http://localhost:8080"
echo ""
npm run dev
