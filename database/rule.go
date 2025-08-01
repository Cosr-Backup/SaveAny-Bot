package database

import "context"

// CreateRule creates a rule in the database
func CreateRule(ctx context.Context, rule *Rule) error {
	return GetDatabase().CreateRule(ctx, rule)
}

// DeleteRule deletes a rule by ID from the database
func DeleteRule(ctx context.Context, ruleID uint) error {
	return GetDatabase().DeleteRule(ctx, ruleID)
}

// UpdateUserApplyRule updates the apply_rule field for a user in the database
func UpdateUserApplyRule(ctx context.Context, chatID int64, applyRule bool) error {
	return GetDatabase().UpdateUserApplyRule(ctx, chatID, applyRule)
}

// GetRulesByUserChatID retrieves rules for a user by chat ID from the database
func GetRulesByUserChatID(ctx context.Context, chatID int64) ([]Rule, error) {
	return GetDatabase().GetRulesByUserChatID(ctx, chatID)
}
