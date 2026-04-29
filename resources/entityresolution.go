package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/entityresolution"
)

// EntityresolutionClient is the interface for the entityresolution SDK client methods.
type EntityresolutionClient interface {
	ListSchemaMappings(ctx context.Context, params *entityresolution.ListSchemaMappingsInput,
		optFns ...func(*entityresolution.Options)) (*entityresolution.ListSchemaMappingsOutput, error)
	DeleteSchemaMapping(ctx context.Context, params *entityresolution.DeleteSchemaMappingInput,
		optFns ...func(*entityresolution.Options)) (*entityresolution.DeleteSchemaMappingOutput, error)
}
