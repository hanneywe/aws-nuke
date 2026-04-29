package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

// CognitoClient is the interface for the Cognito Identity Provider SDK v2 client methods.
type CognitoClient interface {
	ListUserPools(ctx context.Context, params *cognitoidentityprovider.ListUserPoolsInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUserPoolsOutput, error)
	ListResourceServers(ctx context.Context, params *cognitoidentityprovider.ListResourceServersInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListResourceServersOutput, error)
	DeleteResourceServer(ctx context.Context, params *cognitoidentityprovider.DeleteResourceServerInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DeleteResourceServerOutput, error)
	ListUserImportJobs(ctx context.Context, params *cognitoidentityprovider.ListUserImportJobsInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUserImportJobsOutput, error)
	StopUserImportJob(ctx context.Context, params *cognitoidentityprovider.StopUserImportJobInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.StopUserImportJobOutput, error)
	DescribeManagedLoginBrandingByClient(ctx context.Context, params *cognitoidentityprovider.DescribeManagedLoginBrandingByClientInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DescribeManagedLoginBrandingByClientOutput, error)
	DeleteManagedLoginBranding(ctx context.Context, params *cognitoidentityprovider.DeleteManagedLoginBrandingInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.DeleteManagedLoginBrandingOutput, error)
	ListUserPoolClients(ctx context.Context, params *cognitoidentityprovider.ListUserPoolClientsInput,
		optFns ...func(*cognitoidentityprovider.Options)) (*cognitoidentityprovider.ListUserPoolClientsOutput, error)
}
