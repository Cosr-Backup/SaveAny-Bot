package database

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
	"github.com/redis/go-redis/v9"
)

// RedisDatabase Redis数据库实现
type RedisDatabase struct {
	rdb *redis.Client
}

// NewRedisDatabase 创建新的Redis数据库实例
func NewRedisDatabase(ctx context.Context) (Database, error) {
	logger := log.FromContext(ctx)
	
	// 创建Redis客户端选项
	opts := &redis.Options{
		Addr:     config.Cfg.DB.RedisAddr,
		Password: config.Cfg.DB.RedisPassword,
		DB:       config.Cfg.DB.RedisDB,
	}
	
	// 如果配置了Redis ACL认证用户名，则设置
	if config.Cfg.DB.RedisUser != "" {
		opts.Username = config.Cfg.DB.RedisUser
		logger.Debug("Redis ACL username configured", "username", config.Cfg.DB.RedisUser)
	}
	
	// 创建Redis客户端
	rdb := redis.NewClient(opts)

	// 测试连接
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	
	logger.Debug("Redis connected", "ping", pong)
	
	redisDB := &RedisDatabase{
		rdb: rdb,
	}
	
	return redisDB, nil
}

// CreateUser 创建用户
func (r *RedisDatabase) CreateUser(ctx context.Context, chatID int64) error {
	return redisCreateUser(ctx, chatID, r.rdb)
}

// GetAllUsers 获取所有用户
func (r *RedisDatabase) GetAllUsers(ctx context.Context) ([]User, error) {
	return redisGetAllUsers(ctx, r.rdb)
}

// GetUserByChatID 根据聊天ID获取用户
func (r *RedisDatabase) GetUserByChatID(ctx context.Context, chatID int64) (*User, error) {
	return redisGetUserByChatID(ctx, chatID, r.rdb)
}

// UpdateUser 更新用户
func (r *RedisDatabase) UpdateUser(ctx context.Context, user *User) error {
	return redisUpdateUser(ctx, user, r.rdb)
}

// DeleteUser 删除用户
func (r *RedisDatabase) DeleteUser(ctx context.Context, user *User) error {
	return redisDeleteUser(ctx, user, r.rdb)
}

// CreateDirForUser 为用户创建目录
func (r *RedisDatabase) CreateDirForUser(ctx context.Context, userID uint, storageName, path string) error {
	return redisCreateDirForUser(ctx, userID, storageName, path, r.rdb)
}

// GetDirByID 根据ID获取目录
func (r *RedisDatabase) GetDirByID(ctx context.Context, id uint) (*Dir, error) {
	// 对于Redis，我们需要先找到userID，因为没有全局目录索引
	// 这是我们Redis设计的一个限制 - 我们需要搜索用户目录
	// 在生产系统中，您可能需要维护一个全局目录索引
	return nil, fmt.Errorf("GetDirByID not efficiently supported with Redis - use GetUserDirs instead")
}

// GetUserDirs 获取用户的所有目录
func (r *RedisDatabase) GetUserDirs(ctx context.Context, userID uint) ([]Dir, error) {
	return redisGetUserDirs(ctx, userID, r.rdb)
}

// GetUserDirsByChatID 根据聊天ID获取用户的所有目录
func (r *RedisDatabase) GetUserDirsByChatID(ctx context.Context, chatID int64) ([]Dir, error) {
	return redisGetUserDirsByChatID(ctx, chatID, r.rdb)
}

// GetDirsByUserIDAndStorageName 根据用户ID和存储名获取目录
func (r *RedisDatabase) GetDirsByUserIDAndStorageName(ctx context.Context, userID uint, storageName string) ([]Dir, error) {
	return redisGetDirsByUserIDAndStorageName(ctx, userID, storageName, r.rdb)
}

// GetDirsByUserChatIDAndStorageName 根据聊天ID和存储名获取目录
func (r *RedisDatabase) GetDirsByUserChatIDAndStorageName(ctx context.Context, chatID int64, storageName string) ([]Dir, error) {
	return redisGetDirsByUserChatIDAndStorageName(ctx, chatID, storageName, r.rdb)
}

// DeleteDirForUser 删除用户的目录
func (r *RedisDatabase) DeleteDirForUser(ctx context.Context, userID uint, storageName, path string) error {
	return redisDeleteDirForUser(ctx, userID, storageName, path, r.rdb)
}

// DeleteDirByID 根据ID删除目录
func (r *RedisDatabase) DeleteDirByID(ctx context.Context, id uint) error {
	return redisDeleteDirByID(ctx, id, r.rdb)
}

// CreateRule 创建规则
func (r *RedisDatabase) CreateRule(ctx context.Context, rule *Rule) error {
	return redisCreateRule(ctx, rule, r.rdb)
}

// DeleteRule 删除规则
func (r *RedisDatabase) DeleteRule(ctx context.Context, ruleID uint) error {
	return redisDeleteRule(ctx, ruleID, r.rdb)
}

// UpdateUserApplyRule 更新用户应用规则设置
func (r *RedisDatabase) UpdateUserApplyRule(ctx context.Context, chatID int64, applyRule bool) error {
	return redisUpdateUserApplyRule(ctx, chatID, applyRule, r.rdb)
}

// GetRulesByUserChatID 根据聊天ID获取用户规则
func (r *RedisDatabase) GetRulesByUserChatID(ctx context.Context, chatID int64) ([]Rule, error) {
	return redisGetRulesByUserChatID(ctx, chatID, r.rdb)
}