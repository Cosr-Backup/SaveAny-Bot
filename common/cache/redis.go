package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

// InitRedis initializes Redis client with ACL user support for Redis 6.0+ and cloud services
func InitRedis() error {
	if redisClient != nil {
		return fmt.Errorf("redis cache already initialized")
	}

	cfg := config.Cfg.Cache.Redis
	if !cfg.Enable {
		return nil // Redis disabled, skip initialization
	}

	// Build Redis connection options with ACL user support
	opts := &redis.Options{
		Addr:            fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:        cfg.Password,
		DB:              cfg.DB,
		MaxRetries:      cfg.MaxRetries,
		MinIdleConns:    cfg.MinIdleConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		MaxActiveConns:  cfg.MaxActiveConns,
		DialTimeout:     time.Duration(cfg.ConnectTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
	}

	// Redis 6.0+ ACL user support for cloud services (AWS ElastiCache, Azure Cache, etc.)
	if cfg.Username != "" {
		opts.Username = cfg.Username
		log.Infof("Redis: Using ACL user '%s' for Redis 6.0+ authentication", cfg.Username)
	}

	// Create Redis client
	redisClient = redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Info("Redis cache initialized successfully")
	if cfg.Username != "" {
		log.Infof("Redis: Connected with ACL user '%s' to %s:%d (DB: %d)", 
			cfg.Username, cfg.Host, cfg.Port, cfg.DB)
	} else {
		log.Infof("Redis: Connected to %s:%d (DB: %d)", cfg.Host, cfg.Port, cfg.DB)
	}

	return nil
}

// RedisSet stores a value in Redis with TTL
func RedisSet(key string, value any) error {
	if redisClient == nil {
		return fmt.Errorf("redis client not initialized")
	}

	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	ctx := context.Background()
	ttl := time.Duration(config.Cfg.Cache.TTL) * time.Second

	if err := redisClient.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set value in Redis: %w", err)
	}

	return nil
}

// RedisGet retrieves a value from Redis and deserializes it
func RedisGet[T any](key string) (T, bool) {
	var zero T

	if redisClient == nil {
		return zero, false
	}

	ctx := context.Background()
	data, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return zero, false // Key not found
		}
		log.Warnf("Redis get failed for key '%s': %v", key, err)
		return zero, false
	}

	var result T
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		log.Warnf("Failed to unmarshal Redis value for key '%s': %v", key, err)
		return zero, false
	}

	return result, true
}

// RedisClose closes the Redis connection
func RedisClose() error {
	if redisClient != nil {
		return redisClient.Close()
	}
	return nil
}

// IsRedisEnabled returns true if Redis cache is enabled
func IsRedisEnabled() bool {
	return config.Cfg.Cache.Redis.Enable && redisClient != nil
}