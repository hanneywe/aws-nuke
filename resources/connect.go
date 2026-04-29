package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"
)

// ConnectClient is the interface for the Connect SDK client methods.
type ConnectClient interface {
	ListInstances(ctx context.Context, params *connect.ListInstancesInput,
		optFns ...func(*connect.Options)) (*connect.ListInstancesOutput, error)
	DeleteInstance(ctx context.Context, params *connect.DeleteInstanceInput,
		optFns ...func(*connect.Options)) (*connect.DeleteInstanceOutput, error)
	ListSecurityProfiles(ctx context.Context, params *connect.ListSecurityProfilesInput,
		optFns ...func(*connect.Options)) (*connect.ListSecurityProfilesOutput, error)
	DeleteSecurityProfile(ctx context.Context, params *connect.DeleteSecurityProfileInput,
		optFns ...func(*connect.Options)) (*connect.DeleteSecurityProfileOutput, error)
	ListPredefinedAttributes(ctx context.Context, params *connect.ListPredefinedAttributesInput,
		optFns ...func(*connect.Options)) (*connect.ListPredefinedAttributesOutput, error)
	DeletePredefinedAttribute(ctx context.Context, params *connect.DeletePredefinedAttributeInput,
		optFns ...func(*connect.Options)) (*connect.DeletePredefinedAttributeOutput, error)
	ListAgentStatuses(ctx context.Context, params *connect.ListAgentStatusesInput,
		optFns ...func(*connect.Options)) (*connect.ListAgentStatusesOutput, error)
	ListHoursOfOperations(ctx context.Context, params *connect.ListHoursOfOperationsInput,
		optFns ...func(*connect.Options)) (*connect.ListHoursOfOperationsOutput, error)
	DeleteHoursOfOperation(ctx context.Context, params *connect.DeleteHoursOfOperationInput,
		optFns ...func(*connect.Options)) (*connect.DeleteHoursOfOperationOutput, error)
	ListQueues(ctx context.Context, params *connect.ListQueuesInput,
		optFns ...func(*connect.Options)) (*connect.ListQueuesOutput, error)
	DeleteQueue(ctx context.Context, params *connect.DeleteQueueInput,
		optFns ...func(*connect.Options)) (*connect.DeleteQueueOutput, error)
	ListQuickConnects(ctx context.Context, params *connect.ListQuickConnectsInput,
		optFns ...func(*connect.Options)) (*connect.ListQuickConnectsOutput, error)
	DeleteQuickConnect(ctx context.Context, params *connect.DeleteQuickConnectInput,
		optFns ...func(*connect.Options)) (*connect.DeleteQuickConnectOutput, error)
	ListIntegrationAssociations(ctx context.Context, params *connect.ListIntegrationAssociationsInput,
		optFns ...func(*connect.Options)) (*connect.ListIntegrationAssociationsOutput, error)
	ListUseCases(ctx context.Context, params *connect.ListUseCasesInput,
		optFns ...func(*connect.Options)) (*connect.ListUseCasesOutput, error)
	DeleteUseCase(ctx context.Context, params *connect.DeleteUseCaseInput,
		optFns ...func(*connect.Options)) (*connect.DeleteUseCaseOutput, error)
	ListWorkspaces(ctx context.Context, params *connect.ListWorkspacesInput,
		optFns ...func(*connect.Options)) (*connect.ListWorkspacesOutput, error)
	DeleteWorkspace(ctx context.Context, params *connect.DeleteWorkspaceInput,
		optFns ...func(*connect.Options)) (*connect.DeleteWorkspaceOutput, error)

	DeleteIntegrationAssociation(ctx context.Context, params *connect.DeleteIntegrationAssociationInput,
		optFns ...func(*connect.Options)) (*connect.DeleteIntegrationAssociationOutput, error)

	ListRoutingProfiles(ctx context.Context, params *connect.ListRoutingProfilesInput,
		optFns ...func(*connect.Options)) (*connect.ListRoutingProfilesOutput, error)
	DeleteRoutingProfile(ctx context.Context, params *connect.DeleteRoutingProfileInput,
		optFns ...func(*connect.Options)) (*connect.DeleteRoutingProfileOutput, error)
}
