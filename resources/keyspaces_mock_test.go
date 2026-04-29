package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockKeyspacesClient struct {
	mock.Mock
}

func (m *mockKeyspacesClient) ListKeyspaces(ctx context.Context,
	params *keyspaces.ListKeyspacesInput,
	_ ...func(*keyspaces.Options)) (*keyspaces.ListKeyspacesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*keyspaces.ListKeyspacesOutput), args.Error(1)
}

func (m *mockKeyspacesClient) DeleteKeyspace(ctx context.Context,
	params *keyspaces.DeleteKeyspaceInput,
	_ ...func(*keyspaces.Options)) (*keyspaces.DeleteKeyspaceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*keyspaces.DeleteKeyspaceOutput), args.Error(1)
}

func (m *mockKeyspacesClient) ListTypes(ctx context.Context,
	params *keyspaces.ListTypesInput,
	_ ...func(*keyspaces.Options)) (*keyspaces.ListTypesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*keyspaces.ListTypesOutput), args.Error(1)
}

func (m *mockKeyspacesClient) DeleteType(ctx context.Context,
	params *keyspaces.DeleteTypeInput,
	_ ...func(*keyspaces.Options)) (*keyspaces.DeleteTypeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*keyspaces.DeleteTypeOutput), args.Error(1)
}

var testKeyspacesListerOpts = &nuke.ListerOpts{}
