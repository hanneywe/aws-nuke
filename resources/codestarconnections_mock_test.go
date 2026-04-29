package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codestarconnections"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCodeStarConnectionsClient struct {
	mock.Mock
}

func (m *mockCodeStarConnectionsClient) ListHosts(
	ctx context.Context, params *codestarconnections.ListHostsInput,
	_ ...func(*codestarconnections.Options),
) (*codestarconnections.ListHostsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codestarconnections.ListHostsOutput), args.Error(1)
}

func (m *mockCodeStarConnectionsClient) DeleteHost(
	ctx context.Context, params *codestarconnections.DeleteHostInput,
	_ ...func(*codestarconnections.Options),
) (*codestarconnections.DeleteHostOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*codestarconnections.DeleteHostOutput), args.Error(1)
}

var testCodeStarConnectionsListerOpts = &nuke.ListerOpts{}
