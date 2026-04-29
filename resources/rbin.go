package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/rbin"
)

// RbinClient is an interface for the Recycle Bin SDK client methods used by all Rbin resources.
// It enables mock testing of List and Remove operations.
// RbinClient is an interface for the Recycle Bin SDK client methods used by all Rbin resources.
// It enables mock testing of List and Remove operations.
type RbinClient interface {
	// Listing
	ListRules(ctx context.Context, params *rbin.ListRulesInput,
		optFns ...func(*rbin.Options)) (*rbin.ListRulesOutput, error)

	// Deletion
	DeleteRule(ctx context.Context, params *rbin.DeleteRuleInput,
		optFns ...func(*rbin.Options)) (*rbin.DeleteRuleOutput, error)

	// Get rule details
	GetRule(ctx context.Context, params *rbin.GetRuleInput,
		optFns ...func(*rbin.Options)) (*rbin.GetRuleOutput, error)

	// Unlock
	UnlockRule(ctx context.Context, params *rbin.UnlockRuleInput,
		optFns ...func(*rbin.Options)) (*rbin.UnlockRuleOutput, error)
}
