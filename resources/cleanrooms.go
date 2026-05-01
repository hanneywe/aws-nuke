package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"
)

// CleanRoomsClient is an interface for the AWS Clean Rooms SDK client methods used by all Clean Rooms resources.
// It enables mock testing of List and Remove operations.
type CleanRoomsClient interface {
	ListCollaborations(ctx context.Context, params *cleanrooms.ListCollaborationsInput,
		optFns ...func(*cleanrooms.Options)) (*cleanrooms.ListCollaborationsOutput, error)
	DeleteCollaboration(ctx context.Context, params *cleanrooms.DeleteCollaborationInput,
		optFns ...func(*cleanrooms.Options)) (*cleanrooms.DeleteCollaborationOutput, error)

	ListMemberships(ctx context.Context, params *cleanrooms.ListMembershipsInput,
		optFns ...func(*cleanrooms.Options)) (*cleanrooms.ListMembershipsOutput, error)
	DeleteMembership(ctx context.Context, params *cleanrooms.DeleteMembershipInput,
		optFns ...func(*cleanrooms.Options)) (*cleanrooms.DeleteMembershipOutput, error)

	ListPrivacyBudgetTemplates(ctx context.Context, params *cleanrooms.ListPrivacyBudgetTemplatesInput,
		optFns ...func(*cleanrooms.Options)) (*cleanrooms.ListPrivacyBudgetTemplatesOutput, error)
	DeletePrivacyBudgetTemplate(ctx context.Context, params *cleanrooms.DeletePrivacyBudgetTemplateInput,
		optFns ...func(*cleanrooms.Options)) (*cleanrooms.DeletePrivacyBudgetTemplateOutput, error)
}
