package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockGameLiftV2Client struct {
	mock.Mock
}

func (m *mockGameLiftV2Client) ListAliases(
	ctx context.Context, params *gamelift.ListAliasesInput,
	_ ...func(*gamelift.Options),
) (*gamelift.ListAliasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*gamelift.ListAliasesOutput), args.Error(1)
}

func (m *mockGameLiftV2Client) DeleteAlias(
	ctx context.Context, params *gamelift.DeleteAliasInput,
	_ ...func(*gamelift.Options),
) (*gamelift.DeleteAliasOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*gamelift.DeleteAliasOutput), args.Error(1)
}

func (m *mockGameLiftV2Client) ListLocations(
	ctx context.Context, params *gamelift.ListLocationsInput,
	_ ...func(*gamelift.Options),
) (*gamelift.ListLocationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*gamelift.ListLocationsOutput), args.Error(1)
}

func (m *mockGameLiftV2Client) DeleteLocation(
	ctx context.Context, params *gamelift.DeleteLocationInput,
	_ ...func(*gamelift.Options),
) (*gamelift.DeleteLocationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*gamelift.DeleteLocationOutput), args.Error(1)
}

func (m *mockGameLiftV2Client) DescribeVpcPeeringAuthorizations(
	ctx context.Context, params *gamelift.DescribeVpcPeeringAuthorizationsInput,
	_ ...func(*gamelift.Options),
) (*gamelift.DescribeVpcPeeringAuthorizationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*gamelift.DescribeVpcPeeringAuthorizationsOutput), args.Error(1)
}

func (m *mockGameLiftV2Client) DeleteVpcPeeringAuthorization(
	ctx context.Context, params *gamelift.DeleteVpcPeeringAuthorizationInput,
	_ ...func(*gamelift.Options),
) (*gamelift.DeleteVpcPeeringAuthorizationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*gamelift.DeleteVpcPeeringAuthorizationOutput), args.Error(1)
}

var testGameLiftV2ListerOpts = &nuke.ListerOpts{}
