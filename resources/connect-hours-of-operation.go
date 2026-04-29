package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectHoursOfOperationResource = "ConnectHoursOfOperation"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectHoursOfOperationResource,
		Scope:    nuke.Account,
		Resource: &ConnectHoursOfOperation{},
		Lister:   &ConnectHoursOfOperationLister{},
	})
}

type ConnectHoursOfOperationLister struct {
	svc ConnectClient
}

func (l *ConnectHoursOfOperationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = connect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	instancePaginator := connect.NewListInstancesPaginator(svc, &connect.ListInstancesInput{})
	for instancePaginator.HasMorePages() {
		instanceResp, err := instancePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, instance := range instanceResp.InstanceSummaryList {
			hooPaginator := connect.NewListHoursOfOperationsPaginator(svc, &connect.ListHoursOfOperationsInput{
				InstanceId: instance.Id,
			})
			for hooPaginator.HasMorePages() {
				hooResp, err := hooPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, hoo := range hooResp.HoursOfOperationSummaryList {
					resources = append(resources, &ConnectHoursOfOperation{
						svc:        svc,
						InstanceID: instance.Id,
						ID:         hoo.Id,
						Name:       hoo.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectHoursOfOperation struct {
	svc        ConnectClient
	InstanceID *string `property:"name=InstanceId"`
	ID         *string `property:"name=Id"`
	Name       *string
}

func (r *ConnectHoursOfOperation) Filter() error {
	if r.Name != nil && *r.Name == "Basic Hours" {
		return fmt.Errorf("cannot delete default hours of operation")
	}
	return nil
}

func (r *ConnectHoursOfOperation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHoursOfOperation(ctx, &connect.DeleteHoursOfOperationInput{
		InstanceId:         r.InstanceID,
		HoursOfOperationId: r.ID,
	})
	return err
}

func (r *ConnectHoursOfOperation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectHoursOfOperation) String() string {
	return *r.Name
}
