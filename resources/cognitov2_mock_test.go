package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCognitoidentityproviderClient struct {
	mock.Mock
}

func (m *mockCognitoidentityproviderClient) ListUserPools(
	ctx context.Context, params *cognitoidentityprovider.ListUserPoolsInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.ListUserPoolsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.ListUserPoolsOutput), args.Error(1)
}

func (m *mockCognitoidentityproviderClient) ListGroups(
	ctx context.Context, params *cognitoidentityprovider.ListGroupsInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.ListGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.ListGroupsOutput), args.Error(1)
}

func (m *mockCognitoidentityproviderClient) DeleteGroup(
	ctx context.Context, params *cognitoidentityprovider.DeleteGroupInput,
	_ ...func(*cognitoidentityprovider.Options),
) (*cognitoidentityprovider.DeleteGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cognitoidentityprovider.DeleteGroupOutput), args.Error(1)
}

var testCognitoidentityproviderListerOpts = &nuke.ListerOpts{}
