package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockConnectClient struct {
	mock.Mock
}

func (m *mockConnectClient) ListInstances(ctx context.Context,
	params *connect.ListInstancesInput,
	_ ...func(*connect.Options)) (*connect.ListInstancesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListInstancesOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteInstance(ctx context.Context,
	params *connect.DeleteInstanceInput,
	_ ...func(*connect.Options)) (*connect.DeleteInstanceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteInstanceOutput), args.Error(1)
}

func (m *mockConnectClient) ListSecurityProfiles(ctx context.Context,
	params *connect.ListSecurityProfilesInput,
	_ ...func(*connect.Options)) (*connect.ListSecurityProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListSecurityProfilesOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteSecurityProfile(ctx context.Context,
	params *connect.DeleteSecurityProfileInput,
	_ ...func(*connect.Options)) (*connect.DeleteSecurityProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteSecurityProfileOutput), args.Error(1)
}

var testConnectListerOpts = &nuke.ListerOpts{}

func (m *mockConnectClient) ListPredefinedAttributes(ctx context.Context,
	params *connect.ListPredefinedAttributesInput,
	_ ...func(*connect.Options)) (*connect.ListPredefinedAttributesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListPredefinedAttributesOutput), args.Error(1)
}

func (m *mockConnectClient) DeletePredefinedAttribute(ctx context.Context,
	params *connect.DeletePredefinedAttributeInput,
	_ ...func(*connect.Options)) (*connect.DeletePredefinedAttributeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeletePredefinedAttributeOutput), args.Error(1)
}

func (m *mockConnectClient) ListAgentStatuses(ctx context.Context,
	params *connect.ListAgentStatusesInput,
	_ ...func(*connect.Options)) (*connect.ListAgentStatusesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListAgentStatusesOutput), args.Error(1)
}

func (m *mockConnectClient) ListHoursOfOperations(ctx context.Context,
	params *connect.ListHoursOfOperationsInput,
	_ ...func(*connect.Options)) (*connect.ListHoursOfOperationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListHoursOfOperationsOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteHoursOfOperation(ctx context.Context,
	params *connect.DeleteHoursOfOperationInput,
	_ ...func(*connect.Options)) (*connect.DeleteHoursOfOperationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteHoursOfOperationOutput), args.Error(1)
}

func (m *mockConnectClient) ListQueues(ctx context.Context,
	params *connect.ListQueuesInput,
	_ ...func(*connect.Options)) (*connect.ListQueuesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListQueuesOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteQueue(ctx context.Context,
	params *connect.DeleteQueueInput,
	_ ...func(*connect.Options)) (*connect.DeleteQueueOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteQueueOutput), args.Error(1)
}

func (m *mockConnectClient) ListQuickConnects(ctx context.Context,
	params *connect.ListQuickConnectsInput,
	_ ...func(*connect.Options)) (*connect.ListQuickConnectsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListQuickConnectsOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteQuickConnect(ctx context.Context,
	params *connect.DeleteQuickConnectInput,
	_ ...func(*connect.Options)) (*connect.DeleteQuickConnectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteQuickConnectOutput), args.Error(1)
}

func (m *mockConnectClient) ListIntegrationAssociations(ctx context.Context,
	params *connect.ListIntegrationAssociationsInput,
	_ ...func(*connect.Options)) (*connect.ListIntegrationAssociationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListIntegrationAssociationsOutput), args.Error(1)
}

func (m *mockConnectClient) ListUseCases(ctx context.Context,
	params *connect.ListUseCasesInput,
	_ ...func(*connect.Options)) (*connect.ListUseCasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListUseCasesOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteUseCase(ctx context.Context,
	params *connect.DeleteUseCaseInput,
	_ ...func(*connect.Options)) (*connect.DeleteUseCaseOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteUseCaseOutput), args.Error(1)
}

func (m *mockConnectClient) ListWorkspaces(ctx context.Context,
	params *connect.ListWorkspacesInput,
	_ ...func(*connect.Options)) (*connect.ListWorkspacesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListWorkspacesOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteWorkspace(ctx context.Context,
	params *connect.DeleteWorkspaceInput,
	_ ...func(*connect.Options)) (*connect.DeleteWorkspaceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteWorkspaceOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteIntegrationAssociation(
	ctx context.Context, params *connect.DeleteIntegrationAssociationInput,
	_ ...func(*connect.Options),
) (*connect.DeleteIntegrationAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteIntegrationAssociationOutput), args.Error(1)
}

func (m *mockConnectClient) ListRoutingProfiles(
	ctx context.Context, params *connect.ListRoutingProfilesInput,
	_ ...func(*connect.Options),
) (*connect.ListRoutingProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.ListRoutingProfilesOutput), args.Error(1)
}

func (m *mockConnectClient) DeleteRoutingProfile(
	ctx context.Context, params *connect.DeleteRoutingProfileInput,
	_ ...func(*connect.Options),
) (*connect.DeleteRoutingProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connect.DeleteRoutingProfileOutput), args.Error(1)
}
