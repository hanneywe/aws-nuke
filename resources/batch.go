package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/batch"
)

// BatchClient is the interface for the Batch SDK v2 client methods used by sub-resources.
type BatchClient interface {
	ListSchedulingPolicies(ctx context.Context, params *batch.ListSchedulingPoliciesInput,
		optFns ...func(*batch.Options)) (*batch.ListSchedulingPoliciesOutput, error)
	DeleteSchedulingPolicy(ctx context.Context, params *batch.DeleteSchedulingPolicyInput,
		optFns ...func(*batch.Options)) (*batch.DeleteSchedulingPolicyOutput, error)
	ListConsumableResources(ctx context.Context, params *batch.ListConsumableResourcesInput,
		optFns ...func(*batch.Options)) (*batch.ListConsumableResourcesOutput, error)
	DeleteConsumableResource(ctx context.Context, params *batch.DeleteConsumableResourceInput,
		optFns ...func(*batch.Options)) (*batch.DeleteConsumableResourceOutput, error)

	DescribeJobDefinitions(ctx context.Context, params *batch.DescribeJobDefinitionsInput,
		optFns ...func(*batch.Options)) (*batch.DescribeJobDefinitionsOutput, error)
	DeregisterJobDefinition(ctx context.Context, params *batch.DeregisterJobDefinitionInput,
		optFns ...func(*batch.Options)) (*batch.DeregisterJobDefinitionOutput, error)
}
