package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// CognitoidentityproviderClient is the interface for the cognitoidentityprovider SDK client methods.
type CognitoidentityproviderClient interface {
	ListUserPools(ctx context.Context, params *cognitoidentityprovider.ListUserPoolsInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUserPoolsOutput, error)
	ListGroups(ctx context.Context, params *cognitoidentityprovider.ListGroupsInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListGroupsOutput, error)
	DeleteGroup(ctx context.Context, params *cognitoidentityprovider.DeleteGroupInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DeleteGroupOutput, error)
}
