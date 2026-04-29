package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codeconnections"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCodeConnectionsClient struct {
	mock.Mock
}

func (m *mockCodeConnectionsClient) ListRepositoryLinks(
	ctx context.Context, params *codeconnections.ListRepositoryLinksInput,
	_ ...func(*codeconnections.Options),
) (*codeconnections.ListRepositoryLinksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeconnections.ListRepositoryLinksOutput), args.Error(1)
}

func (m *mockCodeConnectionsClient) DeleteRepositoryLink(
	ctx context.Context, params *codeconnections.DeleteRepositoryLinkInput,
	_ ...func(*codeconnections.Options),
) (*codeconnections.DeleteRepositoryLinkOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeconnections.DeleteRepositoryLinkOutput), args.Error(1)
}

func (m *mockCodeConnectionsClient) ListConnections(
	ctx context.Context, params *codeconnections.ListConnectionsInput,
	_ ...func(*codeconnections.Options),
) (*codeconnections.ListConnectionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeconnections.ListConnectionsOutput), args.Error(1)
}

func (m *mockCodeConnectionsClient) DeleteConnection(
	ctx context.Context, params *codeconnections.DeleteConnectionInput,
	_ ...func(*codeconnections.Options),
) (*codeconnections.DeleteConnectionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codeconnections.DeleteConnectionOutput), args.Error(1)
}

var testCodeConnectionsListerOpts = &nuke.ListerOpts{}
