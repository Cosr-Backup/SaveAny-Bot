package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// SQLiteDatabase SQLite数据库实现
type SQLiteDatabase struct {
	db *gorm.DB
}

// NewSQLiteDatabase 创建新的SQLite数据库实例
func NewSQLiteDatabase(ctx context.Context) (Database, error) {
	logger := log.FromContext(ctx)
	
	if err := os.MkdirAll(filepath.Dir(config.Cfg.DB.Path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	
	db, err := gorm.Open(gormlite.Open(config.Cfg.DB.Path), &gorm.Config{
		Logger: glogger.New(logger, glogger.Config{
			Colorful:                  true,
			SlowThreshold:             time.Second * 5,
			LogLevel:                  glogger.Error,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		}),
		PrepareStmt: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	logger.Debug("Database connected")
	if err := db.AutoMigrate(&User{}, &Dir{}, &Rule{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	
	return &SQLiteDatabase{db: db}, nil
}

// CreateUser 创建用户
func (s *SQLiteDatabase) CreateUser(ctx context.Context, chatID int64) error {
	if _, err := s.GetUserByChatID(ctx, chatID); err == nil {
		return nil
	}
	return s.db.WithContext(ctx).Create(&User{ChatID: chatID}).Error
}

// GetAllUsers 获取所有用户
func (s *SQLiteDatabase) GetAllUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := s.db.Preload("Dirs").
		WithContext(ctx).
		Preload("Rules").
		Find(&users).Error
	return users, err
}

// GetUserByChatID 根据聊天ID获取用户
func (s *SQLiteDatabase) GetUserByChatID(ctx context.Context, chatID int64) (*User, error) {
	var user User
	err := s.db.
		Preload("Dirs").
		WithContext(ctx).
		Preload("Rules").
		Where("chat_id = ?", chatID).First(&user).Error
	return &user, err
}

// UpdateUser 更新用户
func (s *SQLiteDatabase) UpdateUser(ctx context.Context, user *User) error {
	if _, err := s.GetUserByChatID(ctx, user.ChatID); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Save(user).Error
}

// DeleteUser 删除用户
func (s *SQLiteDatabase) DeleteUser(ctx context.Context, user *User) error {
	return s.db.WithContext(ctx).Unscoped().Select("Dirs", "Rules").Delete(user).Error
}

// CreateDirForUser 为用户创建目录
func (s *SQLiteDatabase) CreateDirForUser(ctx context.Context, userID uint, storageName, path string) error {
	dir := Dir{
		UserID:      userID,
		StorageName: storageName,
		Path:        path,
	}
	return s.db.WithContext(ctx).Create(&dir).Error
}

// GetDirByID 根据ID获取目录
func (s *SQLiteDatabase) GetDirByID(ctx context.Context, id uint) (*Dir, error) {
	dir := &Dir{}
	err := s.db.WithContext(ctx).First(dir, id).Error
	if err != nil {
		return nil, err
	}
	return dir, err
}

// GetUserDirs 获取用户的所有目录
func (s *SQLiteDatabase) GetUserDirs(ctx context.Context, userID uint) ([]Dir, error) {
	var dirs []Dir
	err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&dirs).Error
	return dirs, err
}

// GetUserDirsByChatID 根据聊天ID获取用户的所有目录
func (s *SQLiteDatabase) GetUserDirsByChatID(ctx context.Context, chatID int64) ([]Dir, error) {
	user, err := s.GetUserByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return s.GetUserDirs(ctx, user.ID)
}

// GetDirsByUserIDAndStorageName 根据用户ID和存储名获取目录
func (s *SQLiteDatabase) GetDirsByUserIDAndStorageName(ctx context.Context, userID uint, storageName string) ([]Dir, error) {
	var dirs []Dir
	err := s.db.WithContext(ctx).Where("user_id = ? AND storage_name = ?", userID, storageName).Find(&dirs).Error
	return dirs, err
}

// GetDirsByUserChatIDAndStorageName 根据聊天ID和存储名获取目录
func (s *SQLiteDatabase) GetDirsByUserChatIDAndStorageName(ctx context.Context, chatID int64, storageName string) ([]Dir, error) {
	user, err := s.GetUserByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	return s.GetDirsByUserIDAndStorageName(ctx, user.ID, storageName)
}

// DeleteDirForUser 删除用户的目录
func (s *SQLiteDatabase) DeleteDirForUser(ctx context.Context, userID uint, storageName, path string) error {
	return s.db.WithContext(ctx).Unscoped().Where("user_id = ? AND storage_name = ? AND path = ?", userID, storageName, path).Delete(&Dir{}).Error
}

// DeleteDirByID 根据ID删除目录
func (s *SQLiteDatabase) DeleteDirByID(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&Dir{}, id).Error
}

// CreateRule 创建规则
func (s *SQLiteDatabase) CreateRule(ctx context.Context, rule *Rule) error {
	return s.db.WithContext(ctx).Create(rule).Error
}

// DeleteRule 删除规则
func (s *SQLiteDatabase) DeleteRule(ctx context.Context, ruleID uint) error {
	return s.db.WithContext(ctx).Unscoped().Delete(&Rule{}, ruleID).Error
}

// UpdateUserApplyRule 更新用户应用规则设置
func (s *SQLiteDatabase) UpdateUserApplyRule(ctx context.Context, chatID int64, applyRule bool) error {
	return s.db.WithContext(ctx).Model(&User{}).Where("chat_id = ?", chatID).Update("apply_rule", applyRule).Error
}

// GetRulesByUserChatID 根据聊天ID获取用户规则
func (s *SQLiteDatabase) GetRulesByUserChatID(ctx context.Context, chatID int64) ([]Rule, error) {
	var rules []Rule
	err := s.db.WithContext(ctx).Where("user_id = (SELECT id FROM users WHERE chat_id = ?)", chatID).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}