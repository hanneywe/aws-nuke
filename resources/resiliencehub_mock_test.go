package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/resiliencehub"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockResilienceHubClient struct {
	mock.Mock
}

func (m *mockResilienceHubClient) ListApps(
	ctx context.Context, params *resiliencehub.ListAppsInput,
	_ ...func(*resiliencehub.Options),
) (*resiliencehub.ListAppsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*resiliencehub.ListAppsOutput), args.Error(1)
}

func (m *mockResilienceHubClient) DeleteApp(
	ctx context.Context, params *resiliencehub.DeleteAppInput,
	_ ...func(*resiliencehub.Options),
) (*resiliencehub.DeleteAppOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*resiliencehub.DeleteAppOutput), args.Error(1)
}

var testResilienceHubListerOpts = &nuke.ListerOpts{}
