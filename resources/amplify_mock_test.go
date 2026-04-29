package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/amplify"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAmplifyClient struct {
	mock.Mock
}

func (m *mockAmplifyClient) ListApps(
	ctx context.Context, params *amplify.ListAppsInput, _ ...func(*amplify.Options),
) (*amplify.ListAppsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.ListAppsOutput), args.Error(1)
}

func (m *mockAmplifyClient) ListBranches(
	ctx context.Context, params *amplify.ListBranchesInput, _ ...func(*amplify.Options),
) (*amplify.ListBranchesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.ListBranchesOutput), args.Error(1)
}

func (m *mockAmplifyClient) DeleteBranch(
	ctx context.Context, params *amplify.DeleteBranchInput, _ ...func(*amplify.Options),
) (*amplify.DeleteBranchOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.DeleteBranchOutput), args.Error(1)
}

func (m *mockAmplifyClient) ListWebhooks(
	ctx context.Context, params *amplify.ListWebhooksInput, _ ...func(*amplify.Options),
) (*amplify.ListWebhooksOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.ListWebhooksOutput), args.Error(1)
}

func (m *mockAmplifyClient) DeleteWebhook(
	ctx context.Context, params *amplify.DeleteWebhookInput, _ ...func(*amplify.Options),
) (*amplify.DeleteWebhookOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.DeleteWebhookOutput), args.Error(1)
}

func (m *mockAmplifyClient) ListBackendEnvironments(
	ctx context.Context, params *amplify.ListBackendEnvironmentsInput,
	_ ...func(*amplify.Options),
) (*amplify.ListBackendEnvironmentsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.ListBackendEnvironmentsOutput), args.Error(1)
}

func (m *mockAmplifyClient) DeleteBackendEnvironment(
	ctx context.Context, params *amplify.DeleteBackendEnvironmentInput,
	_ ...func(*amplify.Options),
) (*amplify.DeleteBackendEnvironmentOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.DeleteBackendEnvironmentOutput), args.Error(1)
}

func (m *mockAmplifyClient) ListDomainAssociations(
	ctx context.Context, params *amplify.ListDomainAssociationsInput, _ ...func(*amplify.Options),
) (*amplify.ListDomainAssociationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.ListDomainAssociationsOutput), args.Error(1)
}

func (m *mockAmplifyClient) DeleteDomainAssociation(
	ctx context.Context, params *amplify.DeleteDomainAssociationInput, _ ...func(*amplify.Options),
) (*amplify.DeleteDomainAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*amplify.DeleteDomainAssociationOutput), args.Error(1)
}

var testAmplifyListerOpts = &nuke.ListerOpts{}
