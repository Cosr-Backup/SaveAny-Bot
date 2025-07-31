# Redis ACL 支持与缓存配置 / Redis ACL Support and Cache Configuration

## 概述 / Overview

本版本新增了对 Redis 6.0+ ACL 用户认证的完整支持，适用于云服务环境和企业级部署场景。同时保持了与原有 Ristretto 内存缓存的向下兼容性。

This version adds comprehensive support for Redis 6.0+ ACL user authentication, suitable for cloud service environments and enterprise deployment scenarios, while maintaining backward compatibility with the original Ristretto in-memory cache.

## 功能特性 / Features

### ✅ Redis 6.0+ ACL 用户支持 / Redis 6.0+ ACL User Support
- 支持云服务 ACL 认证 / Support for cloud service ACL authentication
- 兼容 AWS ElastiCache, Azure Cache for Redis, Google Cloud Memorystore
- 支持自建 Redis 6.0+ 服务器 / Support for self-hosted Redis 6.0+ servers

### ✅ 双缓存策略 / Dual Cache Strategy
- Redis：用于分布式缓存 / Redis: For distributed caching
- SQLite：Telegram session 专用存储 / SQLite: Dedicated storage for Telegram sessions
- 自动降级：Redis 失败时回退到内存缓存 / Auto fallback: Falls back to in-memory cache when Redis fails

### ✅ 环境变量支持 / Environment Variable Support
- Docker/Kubernetes 友好的配置方式 / Docker/Kubernetes friendly configuration
- 运行时配置覆盖 / Runtime configuration override

## 配置说明 / Configuration Guide

### 1. TOML 配置文件 / TOML Configuration File

```toml
[cache]
# Ristretto 内存缓存设置 / Ristretto in-memory cache settings
ttl = 86400          # 缓存过期时间 (秒) / Cache TTL in seconds
num_counters = 100000 # 缓存计数器数量 / Number of cache counters
max_cost = 1000000   # 最大缓存成本 / Maximum cache cost

  [cache.redis]
  enable = false              # 启用 Redis 缓存 / Enable Redis cache
  host = "localhost"          # Redis 服务器地址 / Redis server host
  port = 6379                 # Redis 服务器端口 / Redis server port
  password = ""               # Redis 密码 / Redis password
  redis_user = ""             # Redis 6.0+ ACL 用户名 / Redis 6.0+ ACL username
  db = 0                      # Redis 数据库编号 / Redis database number
  
  # 连接池设置 / Connection Pool Settings
  max_retries = 3             # 最大重试次数 / Maximum retry attempts
  min_idle_conns = 5          # 最小空闲连接数 / Minimum idle connections
  max_idle_conns = 10         # 最大空闲连接数 / Maximum idle connections
  max_active_conns = 100      # 最大活跃连接数 / Maximum active connections
  
  # 超时设置 (秒) / Timeout Settings (seconds)
  connect_timeout = 10        # 连接超时 / Connection timeout
  read_timeout = 5            # 读取超时 / Read timeout
  write_timeout = 5           # 写入超时 / Write timeout
```

### 2. 环境变量配置 / Environment Variable Configuration

```bash
# 基础 Redis 配置 / Basic Redis Configuration
export REDIS_ENABLE=true
export REDIS_HOST=your-redis-host.com
export REDIS_PORT=6379
export REDIS_PASSWORD=your-password

# Redis 6.0+ ACL 用户 / Redis 6.0+ ACL User
export REDIS_USER=your-acl-username

# 数据库选择 / Database Selection
export REDIS_DB=0
```

### 3. Docker 环境变量 / Docker Environment Variables

```bash
docker run -d \
  -e REDIS_ENABLE=true \
  -e REDIS_HOST=redis.example.com \
  -e REDIS_PORT=6379 \
  -e REDIS_USER=saveany-bot \
  -e REDIS_PASSWORD=your-password \
  -v /path/to/config.toml:/app/config.toml \
  -v /path/to/data:/app/data \
  saveany-bot:latest
```

## 云服务配置示例 / Cloud Service Configuration Examples

### AWS ElastiCache for Redis
```toml
[cache.redis]
enable = true
host = "your-cluster.cache.amazonaws.com"
port = 6379
password = "your-auth-token"
redis_user = "your-iam-user"  # IAM 用户或自定义 ACL 用户 / IAM user or custom ACL user
db = 0
```

### Azure Cache for Redis
```toml
[cache.redis]
enable = true
host = "your-cache.redis.cache.windows.net"
port = 6380  # 通常使用 SSL 端口 / Usually uses SSL port
password = "your-access-key"
redis_user = "your-acl-username"
db = 0
```

### Google Cloud Memorystore
```toml
[cache.redis]
enable = true
host = "your-instance-ip"
port = 6379
password = "your-auth-string"
redis_user = "your-acl-user"
db = 0
```

## 测试指南 / Testing Guide

### 1. 测试 Redis 连接 / Test Redis Connection

```bash
# 构建项目 / Build project
go build -o saveany-bot .

# 测试基本连接 / Test basic connection
redis-cli -h localhost -p 6379 -a your-password ping

# 测试 ACL 用户连接 / Test ACL user connection
redis-cli --user your-acl-username --pass your-password -h localhost -p 6379 ping
```

### 2. 功能测试场景 / Functional Test Scenarios

#### 场景 1：纯内存缓存 / Scenario 1: Pure In-Memory Cache
```toml
[cache.redis]
enable = false  # 禁用 Redis / Disable Redis
```

#### 场景 2：Redis 无 ACL / Scenario 2: Redis without ACL
```toml
[cache.redis]
enable = true
host = "localhost"
port = 6379
password = "your-password"
redis_user = ""  # 留空 / Leave empty
```

#### 场景 3：Redis 6.0+ ACL 用户 / Scenario 3: Redis 6.0+ ACL User
```toml
[cache.redis]
enable = true
host = "localhost"
port = 6379
password = "your-password"
redis_user = "saveany-app"  # ACL 用户 / ACL user
```

#### 场景 4：环境变量覆盖 / Scenario 4: Environment Variable Override
```bash
export REDIS_ENABLE=true
export REDIS_HOST=redis.example.com
export REDIS_USER=prod-user
# config.toml 中的设置将被覆盖 / Settings in config.toml will be overridden
```

### 3. 故障转移测试 / Failover Testing

1. **Redis 连接失败测试** / **Redis Connection Failure Test**
   - 配置错误的 Redis 地址 / Configure incorrect Redis address
   - 验证自动回退到内存缓存 / Verify automatic fallback to in-memory cache

2. **ACL 权限测试** / **ACL Permission Test**
   - 使用权限不足的用户 / Use user with insufficient permissions
   - 验证错误处理和日志记录 / Verify error handling and logging

## 迁移指南 / Migration Guide

### 从现有版本升级 / Upgrading from Existing Version

#### 步骤 1：备份配置 / Step 1: Backup Configuration
```bash
cp config.toml config.toml.backup
```

#### 步骤 2：更新配置文件 / Step 2: Update Configuration File
1. 添加 Redis 配置节 / Add Redis configuration section
2. 根据需要配置 ACL 用户 / Configure ACL user as needed
3. 保持原有缓存设置 / Keep existing cache settings

#### 步骤 3：渐进式启用 / Step 3: Gradual Enablement
```toml
# 第一阶段：禁用 Redis，测试兼容性 / Phase 1: Disable Redis, test compatibility
[cache.redis]
enable = false

# 第二阶段：启用 Redis，无 ACL / Phase 2: Enable Redis, no ACL
[cache.redis]
enable = true
redis_user = ""

# 第三阶段：启用 ACL 用户 / Phase 3: Enable ACL user
[cache.redis]
enable = true
redis_user = "your-acl-user"
```

### 兼容性说明 / Compatibility Notes

- ✅ **数据库存储**：SQLite 数据库文件保持不变 / **Database Storage**: SQLite database files remain unchanged
- ✅ **Telegram Session**：Session 数据仍使用本地 SQLite / **Telegram Session**: Session data still uses local SQLite
- ✅ **配置向下兼容**：现有配置无需修改即可运行 / **Backward Compatible Config**: Existing configurations work without modification
- ✅ **环境变量**：新增环境变量支持，不影响现有部署 / **Environment Variables**: New environment variable support doesn't affect existing deployments

## 故障排除 / Troubleshooting

### 常见问题 / Common Issues

#### 1. Redis 连接失败 / Redis Connection Failed
```
[ERROR] Failed to connect to Redis: connection refused
```
**解决方案 / Solution:**
- 检查 Redis 服务状态 / Check Redis service status
- 验证网络连接和防火墙设置 / Verify network connection and firewall settings
- 确认端口配置正确 / Confirm port configuration is correct

#### 2. ACL 认证失败 / ACL Authentication Failed
```
[ERROR] Redis authentication failed: WRONGPASS invalid username-password pair
```
**解决方案 / Solution:**
- 验证 ACL 用户名和密码 / Verify ACL username and password
- 检查用户权限设置 / Check user permission settings
- 确认 Redis 版本支持 ACL (6.0+) / Confirm Redis version supports ACL (6.0+)

#### 3. 环境变量未生效 / Environment Variables Not Taking Effect
```
[INFO] Redis cache disabled, using Ristretto in-memory cache
```
**解决方案 / Solution:**
- 确认环境变量名称正确 / Confirm environment variable names are correct
- 检查变量值格式 / Check variable value format
- 重启应用程序 / Restart application

### 调试技巧 / Debugging Tips

1. **启用详细日志** / **Enable Verbose Logging**
   ```bash
   export LOG_LEVEL=debug
   ./saveany-bot
   ```

2. **测试 Redis 连接** / **Test Redis Connection**
   ```bash
   redis-cli --user your-username --pass your-password -h your-host -p your-port ping
   ```

3. **监控缓存使用** / **Monitor Cache Usage**
   - 查看应用日志中的缓存相关信息 / Check cache-related information in application logs
   - 使用 Redis 监控工具 / Use Redis monitoring tools

## 性能考虑 / Performance Considerations

- **内存缓存**：延迟最低，但不支持分布式 / **In-Memory Cache**: Lowest latency, but no distributed support
- **Redis 缓存**：支持分布式，适合多实例部署 / **Redis Cache**: Supports distribution, suitable for multi-instance deployment
- **混合策略**：关键数据使用 Redis，临时数据使用内存 / **Hybrid Strategy**: Use Redis for critical data, memory for temporary data

---

更多详细信息请参考：[官方文档](https://sabot.unv.app/)
For more detailed information, please refer to: [Official Documentation](https://sabot.unv.app/)