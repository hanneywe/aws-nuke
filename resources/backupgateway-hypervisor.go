package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/backupgateway"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BackupGatewayHypervisorResource = "BackupGatewayHypervisor"

func init() {
	registry.Register(&registry.Registration{
		Name:     BackupGatewayHypervisorResource,
		Scope:    nuke.Account,
		Resource: &BackupGatewayHypervisor{},
		Lister:   &BackupGatewayHypervisorLister{},
	})
}

type BackupGatewayHypervisorLister struct {
	svc BackupGatewayClient
}

func (l *BackupGatewayHypervisorLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = backupgateway.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := backupgateway.NewListHypervisorsPaginator(svc, &backupgateway.ListHypervisorsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Hypervisors {
			resources = append(resources, &BackupGatewayHypervisor{
				svc:           svc,
				HypervisorArn: item.HypervisorArn,
				Name:          item.Name,
			})
		}
	}

	return resources, nil
}

type BackupGatewayHypervisor struct {
	svc           BackupGatewayClient
	HypervisorArn *string
	Name          *string
}

func (r *BackupGatewayHypervisor) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHypervisor(ctx, &backupgateway.DeleteHypervisorInput{
		HypervisorArn: r.HypervisorArn,
	})
	return err
}

func (r *BackupGatewayHypervisor) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BackupGatewayHypervisor) String() string {
	return *r.HypervisorArn
}
