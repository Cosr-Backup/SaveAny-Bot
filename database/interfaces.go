package database

import (
	"context"
)

// Database 定义数据库操作接口
type Database interface {
	// User operations
	CreateUser(ctx context.Context, chatID int64) error
	GetAllUsers(ctx context.Context) ([]User, error)
	GetUserByChatID(ctx context.Context, chatID int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, user *User) error
	
	// Directory operations
	CreateDirForUser(ctx context.Context, userID uint, storageName, path string) error
	GetDirByID(ctx context.Context, id uint) (*Dir, error)
	GetUserDirs(ctx context.Context, userID uint) ([]Dir, error)
	GetUserDirsByChatID(ctx context.Context, chatID int64) ([]Dir, error)
	GetDirsByUserIDAndStorageName(ctx context.Context, userID uint, storageName string) ([]Dir, error)
	GetDirsByUserChatIDAndStorageName(ctx context.Context, chatID int64, storageName string) ([]Dir, error)
	DeleteDirForUser(ctx context.Context, userID uint, storageName, path string) error
	DeleteDirByID(ctx context.Context, id uint) error
	
	// Rule operations
	CreateRule(ctx context.Context, rule *Rule) error
	DeleteRule(ctx context.Context, ruleID uint) error
	UpdateUserApplyRule(ctx context.Context, chatID int64, applyRule bool) error
	GetRulesByUserChatID(ctx context.Context, chatID int64) ([]Rule, error)
}