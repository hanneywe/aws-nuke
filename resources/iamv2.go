package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

// IAMClient is the interface for the iam SDK client methods.
type IAMClient interface {
	ListAccountAliases(ctx context.Context, params *iam.ListAccountAliasesInput,
		optFns ...func(*iam.Options)) (*iam.ListAccountAliasesOutput, error)
	DeleteAccountAlias(ctx context.Context, params *iam.DeleteAccountAliasInput,
		optFns ...func(*iam.Options)) (*iam.DeleteAccountAliasOutput, error)
}
