package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53RecoveryControlConfigControlPanelResource = "Route53RecoveryControlConfigControlPanel"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53RecoveryControlConfigControlPanelResource,
		Scope:    nuke.Account,
		Resource: &Route53RecoveryControlConfigControlPanel{},
		Lister:   &Route53RecoveryControlConfigControlPanelLister{},
		DependsOn: []string{
			Route53RecoveryControlConfigRoutingControlResource,
		},
	})
}

type Route53RecoveryControlConfigControlPanelLister struct {
	svc Route53RecoveryControlConfigClient
}

func (l *Route53RecoveryControlConfigControlPanelLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = route53recoverycontrolconfig.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := route53recoverycontrolconfig.NewListControlPanelsPaginator(svc, &route53recoverycontrolconfig.ListControlPanelsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cp := range resp.ControlPanels {
			resources = append(resources, &Route53RecoveryControlConfigControlPanel{
				svc:                 svc,
				ControlPanelArn:     cp.ControlPanelArn,
				Name:                cp.Name,
				DefaultControlPanel: cp.DefaultControlPanel,
			})
		}
	}
	return resources, nil
}

type Route53RecoveryControlConfigControlPanel struct {
	svc                 Route53RecoveryControlConfigClient
	ControlPanelArn     *string
	Name                *string
	DefaultControlPanel *bool
}

func (r *Route53RecoveryControlConfigControlPanel) Filter() error {
	if r.DefaultControlPanel != nil && *r.DefaultControlPanel {
		return fmt.Errorf("cannot delete default control panel")
	}
	return nil
}

func (r *Route53RecoveryControlConfigControlPanel) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteControlPanel(ctx, &route53recoverycontrolconfig.DeleteControlPanelInput{
		ControlPanelArn: r.ControlPanelArn,
	})
	return err
}

func (r *Route53RecoveryControlConfigControlPanel) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53RecoveryControlConfigControlPanel) String() string {
	return *r.Name
}
