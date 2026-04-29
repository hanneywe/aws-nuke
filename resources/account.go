package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/account"
)

// AccountClient is an interface for the Account SDK client methods used by all Account resources.
// It enables mock testing of List and Remove operations.
type AccountClient interface {
	GetAlternateContact(ctx context.Context, params *account.GetAlternateContactInput,
		optFns ...func(*account.Options)) (*account.GetAlternateContactOutput, error)
	DeleteAlternateContact(ctx context.Context, params *account.DeleteAlternateContactInput,
		optFns ...func(*account.Options)) (*account.DeleteAlternateContactOutput, error)
}
