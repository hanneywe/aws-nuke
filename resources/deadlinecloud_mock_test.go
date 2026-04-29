package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/deadline"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockDeadlineCloudClient struct {
	mock.Mock
}

func (m *mockDeadlineCloudClient) ListFarms(ctx context.Context,
	params *deadline.ListFarmsInput,
	_ ...func(*deadline.Options)) (*deadline.ListFarmsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.ListFarmsOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) DeleteFarm(ctx context.Context,
	params *deadline.DeleteFarmInput,
	_ ...func(*deadline.Options)) (*deadline.DeleteFarmOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.DeleteFarmOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) ListQueues(ctx context.Context,
	params *deadline.ListQueuesInput,
	_ ...func(*deadline.Options)) (*deadline.ListQueuesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.ListQueuesOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) DeleteQueue(ctx context.Context,
	params *deadline.DeleteQueueInput,
	_ ...func(*deadline.Options)) (*deadline.DeleteQueueOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.DeleteQueueOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) ListStorageProfiles(ctx context.Context,
	params *deadline.ListStorageProfilesInput,
	_ ...func(*deadline.Options)) (*deadline.ListStorageProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.ListStorageProfilesOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) DeleteStorageProfile(ctx context.Context,
	params *deadline.DeleteStorageProfileInput,
	_ ...func(*deadline.Options)) (*deadline.DeleteStorageProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.DeleteStorageProfileOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) ListLimits(ctx context.Context,
	params *deadline.ListLimitsInput,
	_ ...func(*deadline.Options)) (*deadline.ListLimitsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.ListLimitsOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) DeleteLimit(ctx context.Context,
	params *deadline.DeleteLimitInput,
	_ ...func(*deadline.Options)) (*deadline.DeleteLimitOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.DeleteLimitOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) ListQueueLimitAssociations(ctx context.Context,
	params *deadline.ListQueueLimitAssociationsInput,
	_ ...func(*deadline.Options)) (*deadline.ListQueueLimitAssociationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.ListQueueLimitAssociationsOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) DeleteQueueLimitAssociation(ctx context.Context,
	params *deadline.DeleteQueueLimitAssociationInput,
	_ ...func(*deadline.Options)) (*deadline.DeleteQueueLimitAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.DeleteQueueLimitAssociationOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) UpdateQueueLimitAssociation(ctx context.Context,
	params *deadline.UpdateQueueLimitAssociationInput,
	_ ...func(*deadline.Options)) (*deadline.UpdateQueueLimitAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.UpdateQueueLimitAssociationOutput), args.Error(1)
}

func (m *mockDeadlineCloudClient) GetQueueLimitAssociation(ctx context.Context,
	params *deadline.GetQueueLimitAssociationInput,
	_ ...func(*deadline.Options)) (*deadline.GetQueueLimitAssociationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*deadline.GetQueueLimitAssociationOutput), args.Error(1)
}

var testDeadlineCloudListerOpts = &nuke.ListerOpts{}
