package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/synthetics"
)

// SyntheticsClient is the interface for the CloudWatch Synthetics SDK client methods.
type SyntheticsClient interface {
	ListGroups(ctx context.Context, params *synthetics.ListGroupsInput,
		optFns ...func(*synthetics.Options)) (*synthetics.ListGroupsOutput, error)
	DeleteGroup(ctx context.Context, params *synthetics.DeleteGroupInput,
		optFns ...func(*synthetics.Options)) (*synthetics.DeleteGroupOutput, error)
}
