package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/schemas"
)

// SchemasClient is the interface for the EventBridge Schemas SDK client methods.
type SchemasClient interface {
	ListRegistries(ctx context.Context, params *schemas.ListRegistriesInput,
		optFns ...func(*schemas.Options)) (*schemas.ListRegistriesOutput, error)
	DeleteRegistry(ctx context.Context, params *schemas.DeleteRegistryInput,
		optFns ...func(*schemas.Options)) (*schemas.DeleteRegistryOutput, error)
	ListDiscoverers(ctx context.Context, params *schemas.ListDiscoverersInput,
		optFns ...func(*schemas.Options)) (*schemas.ListDiscoverersOutput, error)
	DeleteDiscoverer(ctx context.Context, params *schemas.DeleteDiscovererInput,
		optFns ...func(*schemas.Options)) (*schemas.DeleteDiscovererOutput, error)
}
