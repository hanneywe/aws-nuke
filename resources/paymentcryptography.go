package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/paymentcryptography"
)

// PaymentCryptographyClient is an interface for the Payment Cryptography SDK client methods
// used by all Payment Cryptography resources. It enables mock testing of List and Remove operations.
type PaymentCryptographyClient interface {
	// Listing
	ListKeys(ctx context.Context, params *paymentcryptography.ListKeysInput,
		optFns ...func(*paymentcryptography.Options)) (*paymentcryptography.ListKeysOutput, error)
	ListAliases(ctx context.Context, params *paymentcryptography.ListAliasesInput,
		optFns ...func(*paymentcryptography.Options)) (*paymentcryptography.ListAliasesOutput, error)

	// Deletion
	DeleteKey(ctx context.Context, params *paymentcryptography.DeleteKeyInput,
		optFns ...func(*paymentcryptography.Options)) (*paymentcryptography.DeleteKeyOutput, error)
	DeleteAlias(ctx context.Context, params *paymentcryptography.DeleteAliasInput,
		optFns ...func(*paymentcryptography.Options)) (*paymentcryptography.DeleteAliasOutput, error)
}
