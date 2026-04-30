package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockEKSv2Client struct {
	mock.Mock
}

func (m *mockEKSv2Client) ListEksAnywhereSubscriptions(
	ctx context.Context, params *eks.ListEksAnywhereSubscriptionsInput,
	_ ...func(*eks.Options),
) (*eks.ListEksAnywhereSubscriptionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eks.ListEksAnywhereSubscriptionsOutput), args.Error(1)
}

func (m *mockEKSv2Client) DeleteEksAnywhereSubscription(
	ctx context.Context, params *eks.DeleteEksAnywhereSubscriptionInput,
	_ ...func(*eks.Options),
) (*eks.DeleteEksAnywhereSubscriptionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*eks.DeleteEksAnywhereSubscriptionOutput), args.Error(1)
}

var testEKSv2ListerOpts = &nuke.ListerOpts{}
