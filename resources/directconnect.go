package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"
)

// DirectConnectClient is the interface for the Direct Connect SDK client methods.
type DirectConnectClient interface {
	DescribeDirectConnectGateways(ctx context.Context, params *directconnect.DescribeDirectConnectGatewaysInput,
		optFns ...func(*directconnect.Options)) (*directconnect.DescribeDirectConnectGatewaysOutput, error)
	DeleteDirectConnectGateway(ctx context.Context, params *directconnect.DeleteDirectConnectGatewayInput,
		optFns ...func(*directconnect.Options)) (*directconnect.DeleteDirectConnectGatewayOutput, error)
}
