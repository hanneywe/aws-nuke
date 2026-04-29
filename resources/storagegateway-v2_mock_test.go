package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/storagegateway"
)

type mockStorageGatewayV2Client struct {
	mock.Mock
}

func (m *mockStorageGatewayV2Client) ListTapePools(ctx context.Context, params *storagegateway.ListTapePoolsInput,
	_ ...func(*storagegateway.Options)) (*storagegateway.ListTapePoolsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*storagegateway.ListTapePoolsOutput), args.Error(1)
}

func (m *mockStorageGatewayV2Client) DeleteTapePool(ctx context.Context, params *storagegateway.DeleteTapePoolInput,
	_ ...func(*storagegateway.Options)) (*storagegateway.DeleteTapePoolOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*storagegateway.DeleteTapePoolOutput), args.Error(1)
}
