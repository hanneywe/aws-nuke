package resources

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testCleanRoomsListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

type mockCleanRoomsClient struct {
	mock.Mock
}

func (m *mockCleanRoomsClient) ListCollaborations(ctx context.Context, params *cleanrooms.ListCollaborationsInput,
	_ ...func(*cleanrooms.Options)) (*cleanrooms.ListCollaborationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cleanrooms.ListCollaborationsOutput), args.Error(1)
}

func (m *mockCleanRoomsClient) DeleteCollaboration(ctx context.Context, params *cleanrooms.DeleteCollaborationInput,
	_ ...func(*cleanrooms.Options)) (*cleanrooms.DeleteCollaborationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cleanrooms.DeleteCollaborationOutput), args.Error(1)
}

func (m *mockCleanRoomsClient) ListMemberships(ctx context.Context, params *cleanrooms.ListMembershipsInput,
	_ ...func(*cleanrooms.Options)) (*cleanrooms.ListMembershipsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cleanrooms.ListMembershipsOutput), args.Error(1)
}

func (m *mockCleanRoomsClient) DeleteMembership(ctx context.Context, params *cleanrooms.DeleteMembershipInput,
	_ ...func(*cleanrooms.Options)) (*cleanrooms.DeleteMembershipOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cleanrooms.DeleteMembershipOutput), args.Error(1)
}

func (m *mockCleanRoomsClient) ListPrivacyBudgetTemplates(ctx context.Context, params *cleanrooms.ListPrivacyBudgetTemplatesInput,
	_ ...func(*cleanrooms.Options)) (*cleanrooms.ListPrivacyBudgetTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cleanrooms.ListPrivacyBudgetTemplatesOutput), args.Error(1)
}

func (m *mockCleanRoomsClient) DeletePrivacyBudgetTemplate(ctx context.Context, params *cleanrooms.DeletePrivacyBudgetTemplateInput,
	_ ...func(*cleanrooms.Options)) (*cleanrooms.DeletePrivacyBudgetTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cleanrooms.DeletePrivacyBudgetTemplateOutput), args.Error(1)
}
