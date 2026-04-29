package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/networkmonitor"
	networkmonitortypes "github.com/aws/aws-sdk-go-v2/service/networkmonitor/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NetworkMonitorMonitorResource = "NetworkMonitorMonitor"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkMonitorMonitorResource,
		Scope:    nuke.Account,
		Resource: &NetworkMonitorMonitor{},
		Lister:   &NetworkMonitorMonitorLister{},
	})
}

type NetworkMonitorMonitorLister struct {
	svc NetworkMonitorClient
}

func (l *NetworkMonitorMonitorLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = networkmonitor.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := networkmonitor.NewListMonitorsPaginator(svc, &networkmonitor.ListMonitorsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, monitor := range output.Monitors {
			monitorState := string(monitor.State)
			resources = append(resources, &NetworkMonitorMonitor{
				svc:         svc,
				MonitorName: monitor.MonitorName,
				MonitorArn:  monitor.MonitorArn,
				State:       &monitorState,
				Tags:        monitor.Tags,
			})
		}
	}

	return resources, nil
}

type NetworkMonitorMonitor struct {
	svc         NetworkMonitorClient
	MonitorName *string
	MonitorArn  *string
	State       *string
	Tags        map[string]string
}

func (r *NetworkMonitorMonitor) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMonitor(ctx, &networkmonitor.DeleteMonitorInput{
		MonitorName: r.MonitorName,
	})
	return err
}

func (r *NetworkMonitorMonitor) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NetworkMonitorMonitor) String() string {
	return *r.MonitorName
}

func (r *NetworkMonitorMonitor) Filter() error {
	if r.State != nil && *r.State == string(networkmonitortypes.MonitorStateDeleting) {
		return fmt.Errorf("already %s", *r.State)
	}
	return nil
}
