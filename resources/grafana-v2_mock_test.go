package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/grafana"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockGrafanaV2Client struct {
	mock.Mock
}

func (m *mockGrafanaV2Client) ListWorkspaces(
	ctx context.Context, params *grafana.ListWorkspacesInput,
	_ ...func(*grafana.Options),
) (*grafana.ListWorkspacesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*grafana.ListWorkspacesOutput), args.Error(1)
}

func (m *mockGrafanaV2Client) ListWorkspaceServiceAccounts(
	ctx context.Context, params *grafana.ListWorkspaceServiceAccountsInput,
	_ ...func(*grafana.Options),
) (*grafana.ListWorkspaceServiceAccountsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*grafana.ListWorkspaceServiceAccountsOutput), args.Error(1)
}

func (m *mockGrafanaV2Client) ListWorkspaceServiceAccountTokens(
	ctx context.Context, params *grafana.ListWorkspaceServiceAccountTokensInput,
	_ ...func(*grafana.Options),
) (*grafana.ListWorkspaceServiceAccountTokensOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*grafana.ListWorkspaceServiceAccountTokensOutput), args.Error(1)
}

func (m *mockGrafanaV2Client) DeleteWorkspaceServiceAccountToken(
	ctx context.Context, params *grafana.DeleteWorkspaceServiceAccountTokenInput,
	_ ...func(*grafana.Options),
) (*grafana.DeleteWorkspaceServiceAccountTokenOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*grafana.DeleteWorkspaceServiceAccountTokenOutput), args.Error(1)
}

func (m *mockGrafanaV2Client) DeleteWorkspaceServiceAccount(
	ctx context.Context, params *grafana.DeleteWorkspaceServiceAccountInput,
	_ ...func(*grafana.Options),
) (*grafana.DeleteWorkspaceServiceAccountOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*grafana.DeleteWorkspaceServiceAccountOutput), args.Error(1)
}

var testGrafanaV2ListerOpts = &nuke.ListerOpts{}
