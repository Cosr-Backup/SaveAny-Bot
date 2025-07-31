package cache

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/ristretto/v2"
	"github.com/krau/SaveAny-Bot/config"
)

var cache *ristretto.Cache[string, any]

// Init initializes the cache system (Ristretto or Redis based on configuration)
func Init() {
	// Initialize Ristretto cache for fallback or when Redis is disabled
	if cache != nil {
		panic("ristretto cache already initialized")
	}
	c, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: config.Cfg.Cache.NumCounters,
		MaxCost:     config.Cfg.Cache.MaxCost,
		BufferItems: 64,
		OnReject: func(item *ristretto.Item[any]) {
			log.Warnf("Cache item rejected: key=%s, value=%v", item.Key, item.Value)
		},
	})
	if err != nil {
		log.Fatalf("failed to create ristretto cache: %v", err)
	}
	cache = c
	log.Info("Ristretto in-memory cache initialized")

	// Initialize Redis cache if enabled
	if config.Cfg.Cache.Redis.Enable {
		if err := InitRedis(); err != nil {
			log.Warnf("Failed to initialize Redis cache, falling back to Ristretto: %v", err)
		}
	} else {
		log.Info("Redis cache disabled, using Ristretto in-memory cache")
	}
}

// Set stores a value in the cache (Redis if enabled, otherwise Ristretto)
func Set(key string, value any) error {
	// Try Redis first if enabled
	if IsRedisEnabled() {
		if err := RedisSet(key, value); err != nil {
			log.Warnf("Redis set failed for key '%s', falling back to Ristretto: %v", key, err)
		} else {
			return nil // Redis set successful
		}
	}

	// Fallback to Ristretto
	ok := cache.SetWithTTL(key, value, 0, time.Duration(config.Cfg.Cache.TTL)*time.Second)
	if !ok {
		return fmt.Errorf("failed to set value in ristretto cache")
	}
	cache.Wait()
	return nil
}

// Get retrieves a value from the cache (Redis if enabled, otherwise Ristretto)
func Get[T any](key string) (T, bool) {
	var zero T

	// Try Redis first if enabled
	if IsRedisEnabled() {
		if value, ok := RedisGet[T](key); ok {
			return value, true
		}
		// If Redis fails, continue to Ristretto fallback
	}

	// Fallback to Ristretto
	v, ok := cache.Get(key)
	if !ok {
		return zero, false
	}
	vT, ok := v.(T)
	if !ok {
		return zero, false
	}
	return vT, true
}

// Close closes all cache connections
func Close() error {
	var err error
	if IsRedisEnabled() {
		if redisErr := RedisClose(); redisErr != nil {
			err = redisErr
		}
	}
	if cache != nil {
		cache.Close()
	}
	return err
}
