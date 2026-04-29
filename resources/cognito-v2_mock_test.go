package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCognitoClient struct {
	mock.Mock
}

func (m *mockCognitoClient) ListUserPools(
	ctx context.Context, params *cognitoidentityprovider.ListUserPoolsInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.ListUserPoolsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.ListUserPoolsOutput), args.Error(1)
}

func (m *mockCognitoClient) ListResourceServers(
	ctx context.Context, params *cognitoidentityprovider.ListResourceServersInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.ListResourceServersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.ListResourceServersOutput), args.Error(1)
}

func (m *mockCognitoClient) DeleteResourceServer(
	ctx context.Context, params *cognitoidentityprovider.DeleteResourceServerInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.DeleteResourceServerOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.DeleteResourceServerOutput), args.Error(1)
}

var testCognitoListerOpts = &nuke.ListerOpts{}

func (m *mockCognitoClient) ListUserImportJobs(
	ctx context.Context, params *cognitoidentityprovider.ListUserImportJobsInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.ListUserImportJobsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.ListUserImportJobsOutput), args.Error(1)
}

func (m *mockCognitoClient) StopUserImportJob(
	ctx context.Context, params *cognitoidentityprovider.StopUserImportJobInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.StopUserImportJobOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.StopUserImportJobOutput), args.Error(1)
}

func (m *mockCognitoClient) DescribeManagedLoginBrandingByClient(
	ctx context.Context, params *cognitoidentityprovider.DescribeManagedLoginBrandingByClientInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.DescribeManagedLoginBrandingByClientOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.DescribeManagedLoginBrandingByClientOutput), args.Error(1)
}

func (m *mockCognitoClient) DeleteManagedLoginBranding(
	ctx context.Context, params *cognitoidentityprovider.DeleteManagedLoginBrandingInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.DeleteManagedLoginBrandingOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.DeleteManagedLoginBrandingOutput), args.Error(1)
}

func (m *mockCognitoClient) ListUserPoolClients(
	ctx context.Context, params *cognitoidentityprovider.ListUserPoolClientsInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.ListUserPoolClientsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.ListUserPoolClientsOutput), args.Error(1)
}
