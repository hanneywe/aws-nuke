package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCloudTrailClient struct {
	mock.Mock
}

func (m *mockCloudTrailClient) ListDashboards(
	ctx context.Context, params *cloudtrail.ListDashboardsInput,
	_ ...func(*cloudtrail.Options),
) (*cloudtrail.ListDashboardsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudtrail.ListDashboardsOutput), args.Error(1)
}

func (m *mockCloudTrailClient) GetDashboard(
	ctx context.Context, params *cloudtrail.GetDashboardInput,
	_ ...func(*cloudtrail.Options),
) (*cloudtrail.GetDashboardOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudtrail.GetDashboardOutput), args.Error(1)
}

func (m *mockCloudTrailClient) UpdateDashboard(
	ctx context.Context, params *cloudtrail.UpdateDashboardInput,
	_ ...func(*cloudtrail.Options),
) (*cloudtrail.UpdateDashboardOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudtrail.UpdateDashboardOutput), args.Error(1)
}

func (m *mockCloudTrailClient) DeleteDashboard(
	ctx context.Context, params *cloudtrail.DeleteDashboardInput,
	_ ...func(*cloudtrail.Options),
) (*cloudtrail.DeleteDashboardOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudtrail.DeleteDashboardOutput), args.Error(1)
}

var testCloudTrailListerOpts = &nuke.ListerOpts{}
