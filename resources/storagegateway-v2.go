package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/storagegateway"
)

// StorageGatewayV2Client is an interface for the AWS Storage Gateway SDK v2 client methods.
// This is separate from the existing Storage Gateway resources which use SDK v1.
type StorageGatewayV2Client interface {
	ListTapePools(ctx context.Context, params *storagegateway.ListTapePoolsInput,
		optFns ...func(*storagegateway.Options)) (*storagegateway.ListTapePoolsOutput, error)
	DeleteTapePool(ctx context.Context, params *storagegateway.DeleteTapePoolInput,
		optFns ...func(*storagegateway.Options)) (*storagegateway.DeleteTapePoolOutput, error)
}
