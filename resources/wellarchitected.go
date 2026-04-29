package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/wellarchitected"
)

// WellArchitectedClient is an interface for the AWS Well-Architected Tool SDK client methods
// used by all WellArchitected resources. It enables mock testing of List and Remove operations.
type WellArchitectedClient interface {
	ListWorkloads(ctx context.Context, params *wellarchitected.ListWorkloadsInput,
		optFns ...func(*wellarchitected.Options)) (*wellarchitected.ListWorkloadsOutput, error)
	DeleteWorkload(ctx context.Context, params *wellarchitected.DeleteWorkloadInput,
		optFns ...func(*wellarchitected.Options)) (*wellarchitected.DeleteWorkloadOutput, error)
}
