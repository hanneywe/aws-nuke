package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/supportapp"
)

// SupportAppClient is the interface for the AWS Support App SDK client methods.
type SupportAppClient interface {
	GetAccountAlias(ctx context.Context, params *supportapp.GetAccountAliasInput,
		optFns ...func(*supportapp.Options)) (*supportapp.GetAccountAliasOutput, error)
	DeleteAccountAlias(ctx context.Context, params *supportapp.DeleteAccountAliasInput,
		optFns ...func(*supportapp.Options)) (*supportapp.DeleteAccountAliasOutput, error)
}
