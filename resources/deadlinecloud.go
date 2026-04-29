package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/deadline"
)

// DeadlineCloudClient is the interface for the Deadline Cloud SDK client methods.
type DeadlineCloudClient interface {
	ListFarms(ctx context.Context, params *deadline.ListFarmsInput,
		optFns ...func(*deadline.Options)) (*deadline.ListFarmsOutput, error)
	DeleteFarm(ctx context.Context, params *deadline.DeleteFarmInput,
		optFns ...func(*deadline.Options)) (*deadline.DeleteFarmOutput, error)
	ListQueues(ctx context.Context, params *deadline.ListQueuesInput,
		optFns ...func(*deadline.Options)) (*deadline.ListQueuesOutput, error)
	DeleteQueue(ctx context.Context, params *deadline.DeleteQueueInput,
		optFns ...func(*deadline.Options)) (*deadline.DeleteQueueOutput, error)
	ListStorageProfiles(ctx context.Context, params *deadline.ListStorageProfilesInput,
		optFns ...func(*deadline.Options)) (*deadline.ListStorageProfilesOutput, error)
	DeleteStorageProfile(ctx context.Context, params *deadline.DeleteStorageProfileInput,
		optFns ...func(*deadline.Options)) (*deadline.DeleteStorageProfileOutput, error)
	ListLimits(ctx context.Context, params *deadline.ListLimitsInput,
		optFns ...func(*deadline.Options)) (*deadline.ListLimitsOutput, error)
	DeleteLimit(ctx context.Context, params *deadline.DeleteLimitInput,
		optFns ...func(*deadline.Options)) (*deadline.DeleteLimitOutput, error)
	ListQueueLimitAssociations(ctx context.Context, params *deadline.ListQueueLimitAssociationsInput,
		optFns ...func(*deadline.Options)) (*deadline.ListQueueLimitAssociationsOutput, error)
	DeleteQueueLimitAssociation(ctx context.Context, params *deadline.DeleteQueueLimitAssociationInput,
		optFns ...func(*deadline.Options)) (*deadline.DeleteQueueLimitAssociationOutput, error)
	UpdateQueueLimitAssociation(ctx context.Context, params *deadline.UpdateQueueLimitAssociationInput,
		optFns ...func(*deadline.Options)) (*deadline.UpdateQueueLimitAssociationOutput, error)
	GetQueueLimitAssociation(ctx context.Context, params *deadline.GetQueueLimitAssociationInput,
		optFns ...func(*deadline.Options)) (*deadline.GetQueueLimitAssociationOutput, error)
}
