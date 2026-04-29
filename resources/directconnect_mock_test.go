package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/directconnect"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockDirectConnectClient struct {
	mock.Mock
}

func (m *mockDirectConnectClient) DescribeDirectConnectGateways(ctx context.Context,
	params *directconnect.DescribeDirectConnectGatewaysInput,
	_ ...func(*directconnect.Options)) (*directconnect.DescribeDirectConnectGatewaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*directconnect.DescribeDirectConnectGatewaysOutput), args.Error(1)
}

func (m *mockDirectConnectClient) DeleteDirectConnectGateway(ctx context.Context,
	params *directconnect.DeleteDirectConnectGatewayInput,
	_ ...func(*directconnect.Options)) (*directconnect.DeleteDirectConnectGatewayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*directconnect.DeleteDirectConnectGatewayOutput), args.Error(1)
}

var testDirectConnectListerOpts = &nuke.ListerOpts{}
