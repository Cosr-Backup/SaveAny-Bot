package database

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
)

// databaseInstance 保存当前使用的数据库实例
var databaseInstance Database

// InitDatabase 根据配置初始化数据库实例
func InitDatabase(ctx context.Context) error {
	logger := log.FromContext(ctx)
	
	// 检查是否配置了Redis
	if config.Cfg.DB.RedisAddr != "" {
		logger.Info("Initializing Redis database")
		redisDB, err := NewRedisDatabase(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize Redis database: %w", err)
		}
		databaseInstance = redisDB
		logger.Info("Redis database initialized")
	} else {
		// 使用SQLite
		logger.Info("Initializing SQLite database")
		sqliteDB, err := NewSQLiteDatabase(ctx)
		if err != nil {
			return fmt.Errorf("failed to initialize SQLite database: %w", err)
		}
		databaseInstance = sqliteDB
		logger.Info("SQLite database initialized")
	}
	
	// 同步用户数据
	if err := syncUsers(ctx); err != nil {
		return fmt.Errorf("failed to sync users: %w", err)
	}
	
	return nil
}

// GetDatabase 获取数据库实例
func GetDatabase() Database {
	return databaseInstance
}

// syncUsers 同步配置文件中的用户到数据库
func syncUsers(ctx context.Context) error {
	logger := log.FromContext(ctx)
	dbUsers, err := databaseInstance.GetAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	dbUserMap := make(map[int64]User)
	for _, u := range dbUsers {
		dbUserMap[u.ChatID] = u
	}

	cfgUserMap := make(map[int64]struct{})
	for _, u := range config.Cfg.Users {
		cfgUserMap[u.ID] = struct{}{}
	}

	for cfgID := range cfgUserMap {
		if _, exists := dbUserMap[cfgID]; !exists {
			if err := databaseInstance.CreateUser(ctx, cfgID); err != nil {
				return fmt.Errorf("failed to create user %d: %w", cfgID, err)
			}
			logger.Infof("创建用户: %d", cfgID)
		}
	}

	for dbID, dbUser := range dbUserMap {
		if _, exists := cfgUserMap[dbID]; !exists {
			if err := databaseInstance.DeleteUser(ctx, &dbUser); err != nil {
				logger.Warnf("Failed to delete user %d: %v", dbID, err)
			} else {
				logger.Debugf("Deleted user: %d", dbID)
			}
		}
	}

	return nil
}