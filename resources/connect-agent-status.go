package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/connect"
	connecttypes "github.com/aws/aws-sdk-go-v2/service/connect/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectAgentStatusResource = "ConnectAgentStatus"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectAgentStatusResource,
		Scope:    nuke.Account,
		Resource: &ConnectAgentStatus{},
		Lister:   &ConnectAgentStatusLister{},
	})
}

type ConnectAgentStatusLister struct {
	svc ConnectClient
}

func (l *ConnectAgentStatusLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			statusPaginator := connect.NewListAgentStatusesPaginator(svc, &connect.ListAgentStatusesInput{
				InstanceId: instance.Id,
			})
			for statusPaginator.HasMorePages() {
				statusResp, err := statusPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, status := range statusResp.AgentStatusSummaryList {
					resources = append(resources, &ConnectAgentStatus{
						svc:        svc,
						InstanceID: instance.Id,
						ID:         status.Id,
						Name:       status.Name,
						Type:       status.Type,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectAgentStatus struct {
	svc        ConnectClient
	InstanceID *string `property:"name=InstanceId"`
	ID         *string `property:"name=Id"`
	Name       *string
	Type       connecttypes.AgentStatusType
}

func (r *ConnectAgentStatus) Filter() error {
	if r.Type != connecttypes.AgentStatusTypeCustom {
		return fmt.Errorf("cannot delete default agent status of type %s", r.Type)
	}
	return nil
}

func (r *ConnectAgentStatus) Remove(_ context.Context) error {
	return fmt.Errorf("agent statuses cannot be deleted via the Connect API")
}

func (r *ConnectAgentStatus) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectAgentStatus) String() string {
	return *r.Name
}
