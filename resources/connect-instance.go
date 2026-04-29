package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectInstanceResource = "ConnectInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectInstanceResource,
		Scope:    nuke.Account,
		Resource: &ConnectInstance{},
		Lister:   &ConnectInstanceLister{},
		DependsOn: []string{
			ConnectSecurityProfileResource,
		},
	})
}

type ConnectInstanceLister struct {
	svc ConnectClient
}

func (l *ConnectInstanceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := connect.NewListInstancesPaginator(svc, &connect.ListInstancesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.InstanceSummaryList {
			resources = append(resources, &ConnectInstance{
				svc:           svc,
				ID:            item.Id,
				InstanceAlias: item.InstanceAlias,
			})
		}
	}

	return resources, nil
}

type ConnectInstance struct {
	svc           ConnectClient
	ID            *string `property:"name=Id"`
	InstanceAlias *string
}

func (r *ConnectInstance) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteInstance(ctx, &connect.DeleteInstanceInput{
		InstanceId: r.ID,
	})
	return err
}

func (r *ConnectInstance) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectInstance) String() string {
	return *r.ID
}
