package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/connect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConnectWorkspaceResource = "ConnectWorkspace"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConnectWorkspaceResource,
		Scope:    nuke.Account,
		Resource: &ConnectWorkspace{},
		Lister:   &ConnectWorkspaceLister{},
	})
}

type ConnectWorkspaceLister struct {
	svc ConnectClient
}

func (l *ConnectWorkspaceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			wsPaginator := connect.NewListWorkspacesPaginator(svc, &connect.ListWorkspacesInput{
				InstanceId: instance.Id,
			})
			for wsPaginator.HasMorePages() {
				wsResp, err := wsPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, ws := range wsResp.WorkspaceSummaryList {
					resources = append(resources, &ConnectWorkspace{
						svc:        svc,
						InstanceID: instance.Id,
						ID:         ws.Id,
						Name:       ws.Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type ConnectWorkspace struct {
	svc        ConnectClient
	InstanceID *string `property:"name=InstanceId"`
	ID         *string `property:"name=Id"`
	Name       *string
}

func (r *ConnectWorkspace) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkspace(ctx, &connect.DeleteWorkspaceInput{
		InstanceId:  r.InstanceID,
		WorkspaceId: r.ID,
	})
	return err
}

func (r *ConnectWorkspace) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConnectWorkspace) String() string {
	return *r.ID
}
