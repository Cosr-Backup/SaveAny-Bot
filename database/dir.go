package database

import (
	"context"
	"fmt"
)

// CreateDirForUser creates a directory for a user in the database
func CreateDirForUser(ctx context.Context, userID uint, storageName, path string) error {
	return GetDatabase().CreateDirForUser(ctx, userID, storageName, path)
}

// GetDirByID retrieves a directory by ID from the database
func GetDirByID(ctx context.Context, id uint) (*Dir, error) {
	return GetDatabase().GetDirByID(ctx, id)
}

// GetUserDirs retrieves directories for a user from the database
func GetUserDirs(ctx context.Context, userID uint) ([]Dir, error) {
	return GetDatabase().GetUserDirs(ctx, userID)
}

// GetUserDirsByChatID retrieves directories for a user by chat ID from the database
func GetUserDirsByChatID(ctx context.Context, chatID int64) ([]Dir, error) {
	return GetDatabase().GetUserDirsByChatID(ctx, chatID)
}

// GetDirsByUserIDAndStorageName retrieves directories by user ID and storage name from the database
func GetDirsByUserIDAndStorageName(ctx context.Context, userID uint, storageName string) ([]Dir, error) {
	return GetDatabase().GetDirsByUserIDAndStorageName(ctx, userID, storageName)
}

// GetDirsByUserChatIDAndStorageName retrieves directories by user chat ID and storage name from the database
func GetDirsByUserChatIDAndStorageName(ctx context.Context, chatID int64, storageName string) ([]Dir, error) {
	return GetDatabase().GetDirsByUserChatIDAndStorageName(ctx, chatID, storageName)
}

// DeleteDirForUser deletes a directory for a user from the database
func DeleteDirForUser(ctx context.Context, userID uint, storageName, path string) error {
	return GetDatabase().DeleteDirForUser(ctx, userID, storageName, path)
}

// DeleteDirByID deletes a directory by ID from the database
func DeleteDirByID(ctx context.Context, id uint) error {
	return GetDatabase().DeleteDirByID(ctx, id)
}
