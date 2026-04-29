package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectQuickConnectResource = "ConnectQuickConnect"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectQuickConnectResource,
		Scope:    nuke.Account,
		Resource: &ConnectQuickConnect{},
		Lister:   &ConnectQuickConnectLister{},
	})
}

type ConnectQuickConnectLister struct {
	svc ConnectClient
}

func (l *ConnectQuickConnectLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			qcPaginator := connect.NewListQuickConnectsPaginator(svc, &connect.ListQuickConnectsInput{
				InstanceId: instance.Id,
			})
			for qcPaginator.HasMorePages() {
				qcResp, err := qcPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, qc := range qcResp.QuickConnectSummaryList {
					resources = append(resources, &ConnectQuickConnect{
						svc:        svc,
						InstanceID: instance.Id,
						ID:         qc.Id,
						Name:       qc.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectQuickConnect struct {
	svc        ConnectClient
	InstanceID *string `property:"name=InstanceId"`
	ID         *string `property:"name=Id"`
	Name       *string
}

func (r *ConnectQuickConnect) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteQuickConnect(ctx, &connect.DeleteQuickConnectInput{
		InstanceId:     r.InstanceID,
		QuickConnectId: r.ID,
	})
	return err
}

func (r *ConnectQuickConnect) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectQuickConnect) String() string {
	return *r.Name
}
