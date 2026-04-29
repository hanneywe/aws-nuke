package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockMailManagerClient struct {
	mock.Mock
}

func (m *mockMailManagerClient) ListAddonInstances(
	ctx context.Context, params *mailmanager.ListAddonInstancesInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListAddonInstancesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListAddonInstancesOutput), args.Error(1)
}

func (m *mockMailManagerClient) DeleteAddonInstance(
	ctx context.Context, params *mailmanager.DeleteAddonInstanceInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteAddonInstanceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteAddonInstanceOutput), args.Error(1)
}

func (m *mockMailManagerClient) ListAddonSubscriptions(
	ctx context.Context, params *mailmanager.ListAddonSubscriptionsInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListAddonSubscriptionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListAddonSubscriptionsOutput), args.Error(1)
}

func (m *mockMailManagerClient) DeleteAddonSubscription(
	ctx context.Context, params *mailmanager.DeleteAddonSubscriptionInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteAddonSubscriptionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteAddonSubscriptionOutput), args.Error(1)
}

func (m *mockMailManagerClient) ListAddressLists(
	ctx context.Context, params *mailmanager.ListAddressListsInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListAddressListsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListAddressListsOutput), args.Error(1)
}

func (m *mockMailManagerClient) DeleteAddressList(
	ctx context.Context, params *mailmanager.DeleteAddressListInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteAddressListOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteAddressListOutput), args.Error(1)
}

func (m *mockMailManagerClient) ListArchives(
	ctx context.Context, params *mailmanager.ListArchivesInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListArchivesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListArchivesOutput), args.Error(1)
}

func (m *mockMailManagerClient) DeleteArchive(
	ctx context.Context, params *mailmanager.DeleteArchiveInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteArchiveOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteArchiveOutput), args.Error(1)
}

func (m *mockMailManagerClient) ListRelays(
	ctx context.Context, params *mailmanager.ListRelaysInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListRelaysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListRelaysOutput), args.Error(1)
}

func (m *mockMailManagerClient) DeleteRelay(
	ctx context.Context, params *mailmanager.DeleteRelayInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteRelayOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteRelayOutput), args.Error(1)
}

func (m *mockMailManagerClient) ListRuleSets(
	ctx context.Context, params *mailmanager.ListRuleSetsInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.ListRuleSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.ListRuleSetsOutput), args.Error(1)
}

func (m *mockMailManagerClient) DeleteRuleSet(
	ctx context.Context, params *mailmanager.DeleteRuleSetInput,
	_ ...func(*mailmanager.Options),
) (*mailmanager.DeleteRuleSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*mailmanager.DeleteRuleSetOutput), args.Error(1)
}

var testMailManagerListerOpts = &nuke.ListerOpts{}
