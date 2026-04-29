package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/grafana"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GrafanaWorkspaceServiceAccountResource = "GrafanaWorkspaceServiceAccount"

func init() {
	registry.Register(&registry.Registration{
		Name:     GrafanaWorkspaceServiceAccountResource,
		Scope:    nuke.Account,
		Resource: &GrafanaWorkspaceServiceAccount{},
		Lister:   &GrafanaWorkspaceServiceAccountLister{},
	})
}

type GrafanaWorkspaceServiceAccountLister struct {
	svc GrafanaV2Client
}

func (l *GrafanaWorkspaceServiceAccountLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = grafana.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	workspaceParams := &grafana.ListWorkspacesInput{}
	for {
		workspaceResp, err := svc.ListWorkspaces(ctx, workspaceParams)
		if err != nil {
			return nil, err
		}
		for i := range workspaceResp.Workspaces {
			workspace := &workspaceResp.Workspaces[i]
			saParams := &grafana.ListWorkspaceServiceAccountsInput{
				WorkspaceId: workspace.Id,
			}
			for {
				saResp, err := svc.ListWorkspaceServiceAccounts(ctx, saParams)
				if err != nil {
					return nil, err
				}
				for _, sa := range saResp.ServiceAccounts {
					resources = append(resources, &GrafanaWorkspaceServiceAccount{
						svc:              svc,
						WorkspaceID:      workspace.Id,
						ServiceAccountID: sa.Id,
						Name:             sa.Name,
					})
				}
				if saResp.NextToken == nil {
					break
				}
				saParams.NextToken = saResp.NextToken
			}
		}
		if workspaceResp.NextToken == nil {
			break
		}
		workspaceParams.NextToken = workspaceResp.NextToken
	}

	return resources, nil
}

type GrafanaWorkspaceServiceAccount struct {
	svc              GrafanaV2Client
	WorkspaceID      *string `property:"name=WorkspaceId"`
	ServiceAccountID *string `property:"name=ServiceAccountId"`
	Name             *string
}

func (r *GrafanaWorkspaceServiceAccount) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkspaceServiceAccount(ctx, &grafana.DeleteWorkspaceServiceAccountInput{
		WorkspaceId:      r.WorkspaceID,
		ServiceAccountId: r.ServiceAccountID,
	})
	return err
}

func (r *GrafanaWorkspaceServiceAccount) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GrafanaWorkspaceServiceAccount) String() string {
	return *r.Name
}
