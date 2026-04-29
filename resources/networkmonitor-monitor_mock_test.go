package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/networkmonitor"
	networkmonitortypes "github.com/aws/aws-sdk-go-v2/service/networkmonitor/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testNetworkMonitorListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_NetworkMonitorMonitor_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkMonitorClient)

	mockClient.
		On("ListMonitors", mock.Anything, mock.Anything).
		Return(
			&networkmonitor.ListMonitorsOutput{
				Monitors: []networkmonitortypes.MonitorSummary{
					{
						MonitorName: ptr.String("test-monitor"),
						MonitorArn:  ptr.String("arn:aws:networkmonitor:us-east-1:123456789012:monitor/test-monitor"),
						State:       networkmonitortypes.MonitorStateActive,
						Tags: map[string]string{
							"Environment": "test",
						},
					},
				},
			}, nil,
		)

	lister := &NetworkMonitorMonitorLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkMonitorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	monitor := resources[0].(*NetworkMonitorMonitor)
	assertions.Equal("test-monitor", *monitor.MonitorName)
	assertions.Equal("arn:aws:networkmonitor:us-east-1:123456789012:monitor/test-monitor", *monitor.MonitorArn)
	assertions.Equal("ACTIVE", *monitor.State)
	assertions.Equal("test", monitor.Tags["Environment"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkMonitorMonitor_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkMonitorClient)

	mockClient.
		On("ListMonitors", mock.Anything, mock.Anything).
		Return(
			&networkmonitor.ListMonitorsOutput{
				Monitors: []networkmonitortypes.MonitorSummary{},
			}, nil,
		)

	lister := &NetworkMonitorMonitorLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNetworkMonitorListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkMonitorMonitor_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNetworkMonitorClient)

	monitor := &NetworkMonitorMonitor{
		svc:         mockClient,
		MonitorName: ptr.String("test-monitor"),
		MonitorArn:  ptr.String("arn:aws:networkmonitor:us-east-1:123456789012:monitor/test-monitor"),
		State:       ptr.String("ACTIVE"),
	}

	mockClient.
		On(
			"DeleteMonitor",
			mock.Anything,
			&networkmonitor.DeleteMonitorInput{
				MonitorName: monitor.MonitorName,
			},
		).
		Return(&networkmonitor.DeleteMonitorOutput{}, nil)

	err := monitor.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NetworkMonitorMonitor_Properties(t *testing.T) {
	assertions := assert.New(t)

	monitor := NetworkMonitorMonitor{
		MonitorName: ptr.String("test-monitor"),
		MonitorArn:  ptr.String("arn:aws:networkmonitor:us-east-1:123456789012:monitor/test-monitor"),
		State:       ptr.String("ACTIVE"),
		Tags: map[string]string{
			"Environment": "test",
		},
	}

	properties := monitor.Properties()

	assertions.Equal("test-monitor", properties.Get("MonitorName"))
	assertions.Equal("arn:aws:networkmonitor:us-east-1:123456789012:monitor/test-monitor", properties.Get("MonitorArn"))
	assertions.Equal("ACTIVE", properties.Get("State"))
	assertions.Equal("test", properties.Get("tag:Environment"))
}

func Test_Mock_NetworkMonitorMonitor_String(t *testing.T) {
	assertions := assert.New(t)

	monitor := NetworkMonitorMonitor{
		MonitorName: ptr.String("test-monitor"),
		MonitorArn:  ptr.String("arn:aws:networkmonitor:us-east-1:123456789012:monitor/test-monitor"),
	}

	assertions.Equal("test-monitor", monitor.String())
}

func Test_Mock_NetworkMonitorMonitor_Filter(t *testing.T) {
	assertions := assert.New(t)

	// DELETING state should be filtered
	deletingMonitor := NetworkMonitorMonitor{
		MonitorName: ptr.String("deleting-monitor"),
		MonitorArn:  ptr.String("arn:aws:networkmonitor:us-east-1:123456789012:monitor/deleting-monitor"),
		State:       ptr.String("DELETING"),
	}
	err := deletingMonitor.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETING")

	// ACTIVE state should not be filtered
	activeMonitor := NetworkMonitorMonitor{
		MonitorName: ptr.String("active-monitor"),
		MonitorArn:  ptr.String("arn:aws:networkmonitor:us-east-1:123456789012:monitor/active-monitor"),
		State:       ptr.String("ACTIVE"),
	}
	err = activeMonitor.Filter()
	assertions.NoError(err)
}
