package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/amplify"
)

// AmplifyClient is the interface for the Amplify SDK v2 client methods used by sub-resources.
type AmplifyClient interface {
	ListApps(ctx context.Context, params *amplify.ListAppsInput,
		optFns ...func(*amplify.Options)) (*amplify.ListAppsOutput, error)
	ListBranches(ctx context.Context, params *amplify.ListBranchesInput,
		optFns ...func(*amplify.Options)) (*amplify.ListBranchesOutput, error)
	DeleteBranch(ctx context.Context, params *amplify.DeleteBranchInput,
		optFns ...func(*amplify.Options)) (*amplify.DeleteBranchOutput, error)
	ListWebhooks(ctx context.Context, params *amplify.ListWebhooksInput,
		optFns ...func(*amplify.Options)) (*amplify.ListWebhooksOutput, error)
	DeleteWebhook(ctx context.Context, params *amplify.DeleteWebhookInput,
		optFns ...func(*amplify.Options)) (*amplify.DeleteWebhookOutput, error)
	ListBackendEnvironments(ctx context.Context, params *amplify.ListBackendEnvironmentsInput,
		optFns ...func(*amplify.Options)) (*amplify.ListBackendEnvironmentsOutput, error)
	DeleteBackendEnvironment(ctx context.Context, params *amplify.DeleteBackendEnvironmentInput,
		optFns ...func(*amplify.Options)) (*amplify.DeleteBackendEnvironmentOutput, error)
	ListDomainAssociations(ctx context.Context, params *amplify.ListDomainAssociationsInput,
		optFns ...func(*amplify.Options)) (*amplify.ListDomainAssociationsOutput, error)
	DeleteDomainAssociation(ctx context.Context, params *amplify.DeleteDomainAssociationInput,
		optFns ...func(*amplify.Options)) (*amplify.DeleteDomainAssociationOutput, error)
}
