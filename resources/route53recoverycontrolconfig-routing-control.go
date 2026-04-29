package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53RecoveryControlConfigRoutingControlResource = "Route53RecoveryControlConfigRoutingControl"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53RecoveryControlConfigRoutingControlResource,
		Scope:    nuke.Account,
		Resource: &Route53RecoveryControlConfigRoutingControl{},
		Lister:   &Route53RecoveryControlConfigRoutingControlLister{},
	})
}

type Route53RecoveryControlConfigRoutingControlLister struct {
	svc Route53RecoveryControlConfigClient
}

func (l *Route53RecoveryControlConfigRoutingControlLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = route53recoverycontrolconfig.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First list all control panels
	panelPaginator := route53recoverycontrolconfig.NewListControlPanelsPaginator(svc, &route53recoverycontrolconfig.ListControlPanelsInput{})
	for panelPaginator.HasMorePages() {
		panelResp, err := panelPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cp := range panelResp.ControlPanels {
			// Then list routing controls per panel
			rcPaginator := route53recoverycontrolconfig.NewListRoutingControlsPaginator(svc, &route53recoverycontrolconfig.ListRoutingControlsInput{
				ControlPanelArn: cp.ControlPanelArn,
			})
			for rcPaginator.HasMorePages() {
				rcResp, err := rcPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, rc := range rcResp.RoutingControls {
					resources = append(resources, &Route53RecoveryControlConfigRoutingControl{
						svc:               svc,
						RoutingControlArn: rc.RoutingControlArn,
						Name:              rc.Name,
						ControlPanelArn:   rc.ControlPanelArn,
					})
				}
			}
		}
	}
	return resources, nil
}

type Route53RecoveryControlConfigRoutingControl struct {
	svc               Route53RecoveryControlConfigClient
	RoutingControlArn *string
	Name              *string
	ControlPanelArn   *string
}

func (r *Route53RecoveryControlConfigRoutingControl) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRoutingControl(ctx, &route53recoverycontrolconfig.DeleteRoutingControlInput{
		RoutingControlArn: r.RoutingControlArn,
	})
	return err
}

func (r *Route53RecoveryControlConfigRoutingControl) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53RecoveryControlConfigRoutingControl) String() string {
	return *r.Name
}
