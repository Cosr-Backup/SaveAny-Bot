#!/bin/sh

# 确保数据目录存在，防止 SQLite 初始化失败
# Ensure data directories exist to prevent SQLite initialization failures
mkdir -p /app/data /app/cache /app/downloads
chmod 755 /app/data /app/cache /app/downloads

# 支持通过环境变量配置 Redis (便于 Docker/K8s 部署)
# Support Redis configuration via environment variables (for Docker/K8s deployment)
if [ -n "$REDIS_ENABLE" ]; then
    export SAVEANY_CACHE_REDIS_ENABLE="$REDIS_ENABLE"
fi
if [ -n "$REDIS_HOST" ]; then
    export SAVEANY_CACHE_REDIS_HOST="$REDIS_HOST"
fi
if [ -n "$REDIS_PORT" ]; then
    export SAVEANY_CACHE_REDIS_PORT="$REDIS_PORT"
fi
if [ -n "$REDIS_PASSWORD" ]; then
    export SAVEANY_CACHE_REDIS_PASSWORD="$REDIS_PASSWORD"
fi
if [ -n "$REDIS_USER" ]; then
    export SAVEANY_CACHE_REDIS_REDIS_USER="$REDIS_USER"
fi
if [ -n "$REDIS_DB" ]; then
    export SAVEANY_CACHE_REDIS_DB="$REDIS_DB"
fi

# 下载配置文件 / Download configuration file
if [ -n "$CONFIG_URL" ]; then
    echo "[INFO] Downloading config from $CONFIG_URL"
    if curl -sSLo /app/config.toml "$CONFIG_URL"; then
        echo "[INFO] Configuration downloaded successfully"
    else
        echo "[ERROR] Failed to download config from $CONFIG_URL"
        exit 1
    fi
fi

# 检查配置文件 / Check configuration file
if [ ! -f /app/config.toml ]; then
    echo "[ERROR] Missing config.toml: 请通过挂载或 CONFIG_URL 提供配置文件"
    echo "[ERROR] Missing config.toml: Please provide configuration via mount or CONFIG_URL"
    exit 1
fi

# 显示环境信息 / Display environment info
echo "[INFO] Starting SaveAny-Bot..."
if [ -n "$REDIS_ENABLE" ] && [ "$REDIS_ENABLE" = "true" ]; then
    echo "[INFO] Redis cache enabled via environment variables"
    if [ -n "$REDIS_USER" ]; then
        echo "[INFO] Using Redis ACL user: $REDIS_USER"
    fi
fi
    
exec /app/saveany-bot