package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/outposts"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockOutpostsClient struct {
	mock.Mock
}

func (m *mockOutpostsClient) ListSites(
	ctx context.Context, params *outposts.ListSitesInput,
	_ ...func(*outposts.Options),
) (*outposts.ListSitesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*outposts.ListSitesOutput), args.Error(1)
}

func (m *mockOutpostsClient) DeleteSite(
	ctx context.Context, params *outposts.DeleteSiteInput,
	_ ...func(*outposts.Options),
) (*outposts.DeleteSiteOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*outposts.DeleteSiteOutput), args.Error(1)
}

var testOutpostsListerOpts = &nuke.ListerOpts{}
