package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"
)

// GameLiftV2Client is the interface for the GameLift SDK v2 client methods.
type GameLiftV2Client interface {
	ListAliases(ctx context.Context, params *gamelift.ListAliasesInput,
		optFns ...func(*gamelift.Options)) (*gamelift.ListAliasesOutput, error)
	DeleteAlias(ctx context.Context, params *gamelift.DeleteAliasInput,
		optFns ...func(*gamelift.Options)) (*gamelift.DeleteAliasOutput, error)
	ListLocations(ctx context.Context, params *gamelift.ListLocationsInput,
		optFns ...func(*gamelift.Options)) (*gamelift.ListLocationsOutput, error)
	DeleteLocation(ctx context.Context, params *gamelift.DeleteLocationInput,
		optFns ...func(*gamelift.Options)) (*gamelift.DeleteLocationOutput, error)
	DescribeVpcPeeringAuthorizations(ctx context.Context, params *gamelift.DescribeVpcPeeringAuthorizationsInput,
		optFns ...func(*gamelift.Options)) (*gamelift.DescribeVpcPeeringAuthorizationsOutput, error)
	DeleteVpcPeeringAuthorization(ctx context.Context, params *gamelift.DeleteVpcPeeringAuthorizationInput,
		optFns ...func(*gamelift.Options)) (*gamelift.DeleteVpcPeeringAuthorizationOutput, error)
}
