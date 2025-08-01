package database

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

var db *gorm.DB
var useRedis bool // Flag to determine whether to use Redis or SQLite

// DatabaseInitializer 定义数据库初始化器接口
type DatabaseInitializer interface {
	Init(ctx context.Context) error
}

// SQLiteInitializer 实现SQLite数据库初始化
type SQLiteInitializer struct {
	logger *log.Logger
}

// RedisInitializer 实现Redis数据库初始化
type RedisInitializer struct {
	logger *log.Logger
}

// InitDatabase 使用工厂模式初始化数据库
func InitDatabase(ctx context.Context) error {
	logger := log.FromContext(ctx)
	
	// 工厂模式创建数据库初始化器
	var initializer DatabaseInitializer
	
	// Check if Redis is configured
	if config.Cfg.DB.RedisAddr != "" {
		initializer = &RedisInitializer{logger: logger}
		useRedis = true
	} else {
		initializer = &SQLiteInitializer{logger: logger}
		useRedis = false
		logger.Debug("Redis not configured, using SQLite database")
	}
	
	// 初始化数据库
	if err := initializer.Init(ctx); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	
	// Sync users to database if needed
	if err := syncUsers(ctx); err != nil {
		return fmt.Errorf("failed to sync users: %w", err)
	}
	
	return nil
}

// Init initializes the database (SQLite or Redis based on configuration)
func Init(ctx context.Context) {
	logger := log.FromContext(ctx)
	
	// Initialize database using factory pattern
	if err := InitDatabase(ctx); err != nil {
		logger.Fatal("Failed to initialize database: ", err)
	}
	
	logger.Info("Database initialized successfully")
}

// Init 初始化SQLite数据库
func (i *SQLiteInitializer) Init(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Dir(config.Cfg.DB.Path), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	
	var err error
	db, err = gorm.Open(gormlite.Open(config.Cfg.DB.Path), &gorm.Config{
		Logger: glogger.New(i.logger, glogger.Config{
			Colorful:                  true,
			SlowThreshold:             time.Second * 5,
			LogLevel:                  glogger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		}),
		PrepareStmt: true,
	})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	i.logger.Debug("Database connected")
	if err := db.AutoMigrate(&User{}, &Dir{}, &Rule{}); err != nil {
		return fmt.Errorf("迁移数据库失败, 如果您从旧版本升级, 建议手动删除数据库文件后重试: %w", err)
	}
	
	return nil
}

// Init 初始化Redis数据库
func (i *RedisInitializer) Init(ctx context.Context) error {
	if err := initRedis(ctx); err != nil {
		return fmt.Errorf("failed to initialize redis: %w", err)
	}
	i.logger.Debug("Redis database migrated")
	return nil
}
}

// syncUsers synchronizes users between configuration and database (works with both SQLite and Redis)
func syncUsers(ctx context.Context) error {
	logger := log.FromContext(ctx)
	dbUsers, err := GetAllUsers(ctx)
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
			if err := CreateUser(ctx, cfgID); err != nil {
				return fmt.Errorf("failed to create user %d: %w", cfgID, err)
			}
			logger.Infof("创建用户: %d", cfgID)
		}
	}

	for dbID, dbUser := range dbUserMap {
		if _, exists := cfgUserMap[dbID]; !exists {
			if err := DeleteUser(ctx, &dbUser); err != nil {
				return fmt.Errorf("failed to delete user %d: %w", dbID, err)
			}
			logger.Infof("删除用户: %d", dbID)
		}
	}

	return nil
}
