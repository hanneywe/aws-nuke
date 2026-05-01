package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
)

// BatchComputeEnvironmentClient is an interface for the Batch SDK client methods used by
// BatchComputeEnvironment and BatchComputeEnvironmentState resources.
type BatchComputeEnvironmentClient interface {
	DescribeComputeEnvironments(ctx context.Context, params *batch.DescribeComputeEnvironmentsInput,
		optFns ...func(*batch.Options)) (*batch.DescribeComputeEnvironmentsOutput, error)
	DeleteComputeEnvironment(ctx context.Context, params *batch.DeleteComputeEnvironmentInput,
		optFns ...func(*batch.Options)) (*batch.DeleteComputeEnvironmentOutput, error)
	UpdateComputeEnvironment(ctx context.Context, params *batch.UpdateComputeEnvironmentInput,
		optFns ...func(*batch.Options)) (*batch.UpdateComputeEnvironmentOutput, error)
}

// BatchCEStatusDeleted is the status string for a deleted compute environment.
const BatchCEStatusDeleted = batchtypes.CEStatusDeleted
