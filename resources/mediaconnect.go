package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mediaconnect"
)

// MediaConnectClient is an interface for the MediaConnect SDK client methods used by all MediaConnect resources.
// It enables mock testing of List and Remove operations.
type MediaConnectClient interface {
	// Listing
	ListGateways(ctx context.Context, params *mediaconnect.ListGatewaysInput,
		optFns ...func(*mediaconnect.Options)) (*mediaconnect.ListGatewaysOutput, error)

	// Deletion
	DeleteGateway(ctx context.Context, params *mediaconnect.DeleteGatewayInput,
		optFns ...func(*mediaconnect.Options)) (*mediaconnect.DeleteGatewayOutput, error)
}
