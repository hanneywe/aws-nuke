package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53RecoveryReadinessCellResource = "Route53RecoveryReadinessCell"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53RecoveryReadinessCellResource,
		Scope:    nuke.Account,
		Resource: &Route53RecoveryReadinessCell{},
		Lister:   &Route53RecoveryReadinessCellLister{},
	})
}

type Route53RecoveryReadinessCellLister struct {
	svc Route53RecoveryReadinessClient
}

func (l *Route53RecoveryReadinessCellLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = route53recoveryreadiness.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := route53recoveryreadiness.NewListCellsPaginator(svc, &route53recoveryreadiness.ListCellsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, cell := range resp.Cells {
			resources = append(resources, &Route53RecoveryReadinessCell{
				svc:      svc,
				CellName: cell.CellName,
				CellArn:  cell.CellArn,
			})
		}
	}
	return resources, nil
}

type Route53RecoveryReadinessCell struct {
	svc      Route53RecoveryReadinessClient
	CellName *string
	CellArn  *string
}

func (r *Route53RecoveryReadinessCell) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCell(ctx, &route53recoveryreadiness.DeleteCellInput{
		CellName: r.CellName,
	})
	return err
}

func (r *Route53RecoveryReadinessCell) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53RecoveryReadinessCell) String() string {
	return *r.CellName
}
