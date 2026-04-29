package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/backupgateway"
)

type mockBackupGatewayClient struct {
	mock.Mock
}

func (m *mockBackupGatewayClient) ListHypervisors(ctx context.Context, params *backupgateway.ListHypervisorsInput,
	_ ...func(*backupgateway.Options)) (*backupgateway.ListHypervisorsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*backupgateway.ListHypervisorsOutput), args.Error(1)
}

func (m *mockBackupGatewayClient) DeleteHypervisor(ctx context.Context, params *backupgateway.DeleteHypervisorInput,
	_ ...func(*backupgateway.Options)) (*backupgateway.DeleteHypervisorOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*backupgateway.DeleteHypervisorOutput), args.Error(1)
}
