package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/keyspaces"
)

// KeyspacesClient is the interface for the Keyspaces SDK client methods.
type KeyspacesClient interface {
	ListKeyspaces(ctx context.Context, params *keyspaces.ListKeyspacesInput,
		optFns ...func(*keyspaces.Options)) (*keyspaces.ListKeyspacesOutput, error)
	DeleteKeyspace(ctx context.Context, params *keyspaces.DeleteKeyspaceInput,
		optFns ...func(*keyspaces.Options)) (*keyspaces.DeleteKeyspaceOutput, error)
	ListTypes(ctx context.Context, params *keyspaces.ListTypesInput,
		optFns ...func(*keyspaces.Options)) (*keyspaces.ListTypesOutput, error)
	DeleteType(ctx context.Context, params *keyspaces.DeleteTypeInput,
		optFns ...func(*keyspaces.Options)) (*keyspaces.DeleteTypeOutput, error)
}
