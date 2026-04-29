package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/networkmonitor"
)

// NetworkMonitorClient is an interface for the AWS CloudWatch Network Monitor SDK client methods
// used by all NetworkMonitor resources. It enables mock testing of List and Remove operations.
type NetworkMonitorClient interface {
	ListMonitors(ctx context.Context, params *networkmonitor.ListMonitorsInput,
		optFns ...func(*networkmonitor.Options)) (*networkmonitor.ListMonitorsOutput, error)
	DeleteMonitor(ctx context.Context, params *networkmonitor.DeleteMonitorInput,
		optFns ...func(*networkmonitor.Options)) (*networkmonitor.DeleteMonitorOutput, error)
}
