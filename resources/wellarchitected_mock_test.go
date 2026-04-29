package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
)

type mockWellArchitectedClient struct {
	mock.Mock
}

func (m *mockWellArchitectedClient) ListWorkloads(ctx context.Context, params *wellarchitected.ListWorkloadsInput,
	_ ...func(*wellarchitected.Options)) (*wellarchitected.ListWorkloadsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*wellarchitected.ListWorkloadsOutput), args.Error(1)
}

func (m *mockWellArchitectedClient) DeleteWorkload(ctx context.Context, params *wellarchitected.DeleteWorkloadInput,
	_ ...func(*wellarchitected.Options)) (*wellarchitected.DeleteWorkloadOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*wellarchitected.DeleteWorkloadOutput), args.Error(1)
}
