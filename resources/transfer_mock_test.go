package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/transfer"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testTransferListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockTransferClient struct {
	mock.Mock
}

func (m *mockTransferClient) ListProfiles(ctx context.Context, params *transfer.ListProfilesInput,
	_ ...func(*transfer.Options)) (*transfer.ListProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.ListProfilesOutput), args.Error(1)
}

func (m *mockTransferClient) DeleteProfile(ctx context.Context, params *transfer.DeleteProfileInput,
	_ ...func(*transfer.Options)) (*transfer.DeleteProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.DeleteProfileOutput), args.Error(1)
}

func (m *mockTransferClient) ListConnectors(
	ctx context.Context, params *transfer.ListConnectorsInput,
	_ ...func(*transfer.Options),
) (*transfer.ListConnectorsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.ListConnectorsOutput), args.Error(1)
}

func (m *mockTransferClient) DeleteConnector(
	ctx context.Context, params *transfer.DeleteConnectorInput,
	_ ...func(*transfer.Options),
) (*transfer.DeleteConnectorOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.DeleteConnectorOutput), args.Error(1)
}

func (m *mockTransferClient) ListWorkflows(
	ctx context.Context, params *transfer.ListWorkflowsInput,
	_ ...func(*transfer.Options),
) (*transfer.ListWorkflowsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.ListWorkflowsOutput), args.Error(1)
}

func (m *mockTransferClient) DeleteWorkflow(
	ctx context.Context, params *transfer.DeleteWorkflowInput,
	_ ...func(*transfer.Options),
) (*transfer.DeleteWorkflowOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.DeleteWorkflowOutput), args.Error(1)
}

func (m *mockTransferClient) ListServers(
	ctx context.Context, params *transfer.ListServersInput,
	_ ...func(*transfer.Options),
) (*transfer.ListServersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.ListServersOutput), args.Error(1)
}

func (m *mockTransferClient) ListUsers(
	ctx context.Context, params *transfer.ListUsersInput,
	_ ...func(*transfer.Options),
) (*transfer.ListUsersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.ListUsersOutput), args.Error(1)
}

func (m *mockTransferClient) DescribeUser(
	ctx context.Context, params *transfer.DescribeUserInput,
	_ ...func(*transfer.Options),
) (*transfer.DescribeUserOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.DescribeUserOutput), args.Error(1)
}

func (m *mockTransferClient) DeleteSshPublicKey(
	ctx context.Context, params *transfer.DeleteSshPublicKeyInput,
	_ ...func(*transfer.Options),
) (*transfer.DeleteSshPublicKeyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*transfer.DeleteSshPublicKeyOutput), args.Error(1)
}
