package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
)

// AppRegistryClient is the interface for the ServiceCatalogAppRegistry SDK v2 client methods.
type AppRegistryClient interface {
	ListAttributeGroups(ctx context.Context, params *servicecatalogappregistry.ListAttributeGroupsInput,
		optFns ...func(*servicecatalogappregistry.Options)) (*servicecatalogappregistry.ListAttributeGroupsOutput, error)
	DeleteAttributeGroup(ctx context.Context, params *servicecatalogappregistry.DeleteAttributeGroupInput,
		optFns ...func(*servicecatalogappregistry.Options)) (*servicecatalogappregistry.DeleteAttributeGroupOutput, error)
}
