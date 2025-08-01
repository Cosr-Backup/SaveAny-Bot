package database

import "context"

// CreateUser creates a new user in the database
func CreateUser(ctx context.Context, chatID int64) error {
	return GetDatabase().CreateUser(ctx, chatID)
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(ctx context.Context) ([]User, error) {
	return GetDatabase().GetAllUsers(ctx)
}

// GetUserByChatID retrieves a user by chat ID from the database
func GetUserByChatID(ctx context.Context, chatID int64) (*User, error) {
	return GetDatabase().GetUserByChatID(ctx, chatID)
}

// UpdateUser updates a user in the database
func UpdateUser(ctx context.Context, user *User) error {
	return GetDatabase().UpdateUser(ctx, user)
}

// DeleteUser deletes a user from the database
func DeleteUser(ctx context.Context, user *User) error {
	return GetDatabase().DeleteUser(ctx, user)
}
