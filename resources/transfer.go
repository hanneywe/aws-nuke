package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"
)

// TransferClient is an interface for the Transfer SDK client methods used by all Transfer resources.
// It enables mock testing of List and Remove operations.
type TransferClient interface {
	// Listing
	ListProfiles(ctx context.Context, params *transfer.ListProfilesInput,
		optFns ...func(*transfer.Options)) (*transfer.ListProfilesOutput, error)

	// Deletion
	DeleteProfile(ctx context.Context, params *transfer.DeleteProfileInput,
		optFns ...func(*transfer.Options)) (*transfer.DeleteProfileOutput, error)

	ListConnectors(ctx context.Context, params *transfer.ListConnectorsInput,
		optFns ...func(*transfer.Options)) (*transfer.ListConnectorsOutput, error)
	DeleteConnector(ctx context.Context, params *transfer.DeleteConnectorInput,
		optFns ...func(*transfer.Options)) (*transfer.DeleteConnectorOutput, error)

	ListWorkflows(ctx context.Context, params *transfer.ListWorkflowsInput,
		optFns ...func(*transfer.Options)) (*transfer.ListWorkflowsOutput, error)
	DeleteWorkflow(ctx context.Context, params *transfer.DeleteWorkflowInput,
		optFns ...func(*transfer.Options)) (*transfer.DeleteWorkflowOutput, error)

	ListServers(ctx context.Context, params *transfer.ListServersInput,
		optFns ...func(*transfer.Options)) (*transfer.ListServersOutput, error)
	ListUsers(ctx context.Context, params *transfer.ListUsersInput,
		optFns ...func(*transfer.Options)) (*transfer.ListUsersOutput, error)
	DescribeUser(ctx context.Context, params *transfer.DescribeUserInput,
		optFns ...func(*transfer.Options)) (*transfer.DescribeUserOutput, error)
	DeleteSshPublicKey(ctx context.Context, params *transfer.DeleteSshPublicKeyInput,
		optFns ...func(*transfer.Options)) (*transfer.DeleteSshPublicKeyOutput, error)
}
