package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/organizations"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockOrganizationsClient struct {
	mock.Mock
}

func (m *mockOrganizationsClient) ListAccounts(
	ctx context.Context, params *organizations.ListAccountsInput,
	_ ...func(*organizations.Options),
) (*organizations.ListAccountsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*organizations.ListAccountsOutput), args.Error(1)
}

func (m *mockOrganizationsClient) CloseAccount(
	ctx context.Context, params *organizations.CloseAccountInput,
	_ ...func(*organizations.Options),
) (*organizations.CloseAccountOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*organizations.CloseAccountOutput), args.Error(1)
}

func (m *mockOrganizationsClient) DescribeOrganization(
	ctx context.Context, params *organizations.DescribeOrganizationInput,
	_ ...func(*organizations.Options),
) (*organizations.DescribeOrganizationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*organizations.DescribeOrganizationOutput), args.Error(1)
}

var testOrganizationsListerOpts = &nuke.ListerOpts{}
