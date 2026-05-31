#!/bin/bash

set -e

echo "🐇 Starting RabbitMQ with Docker Compose..."

COMPOSE_DIR="$(dirname "$0")"
cd "$COMPOSE_DIR"

# 检查 Docker
if ! docker info >/dev/null 2>&1; then
  echo "❌ Docker / OrbStack is not running. Please start it first."
  exit 1
fi

# 创建目录（防止 noowners 权限问题）
echo "📁 Ensuring data directories..."
mkdir -p /Volumes/HIKSEMI/mq/rabbitmq/{mnesia,log}

# 停止旧容器（如果存在）
docker compose down || true

# 启动
docker compose up -d

echo ""
echo "✅ RabbitMQ started!"
echo "🌐 Management UI: http://localhost:15672"
echo "👤 Username: admin"
echo "🔑 Password: admin123"
echo ""
echo "📜 Logs:"
docker compose logs --tail 30
