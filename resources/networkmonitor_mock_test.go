package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/networkmonitor"
)

type mockNetworkMonitorClient struct {
	mock.Mock
}

func (m *mockNetworkMonitorClient) ListMonitors(ctx context.Context, params *networkmonitor.ListMonitorsInput,
	_ ...func(*networkmonitor.Options)) (*networkmonitor.ListMonitorsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmonitor.ListMonitorsOutput), args.Error(1)
}

func (m *mockNetworkMonitorClient) DeleteMonitor(ctx context.Context, params *networkmonitor.DeleteMonitorInput,
	_ ...func(*networkmonitor.Options)) (*networkmonitor.DeleteMonitorOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*networkmonitor.DeleteMonitorOutput), args.Error(1)
}
