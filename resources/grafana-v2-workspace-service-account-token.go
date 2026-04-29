package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/grafana"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GrafanaWorkspaceServiceAccountTokenResource = "GrafanaWorkspaceServiceAccountToken"

func init() {
	registry.Register(&registry.Registration{
		Name:     GrafanaWorkspaceServiceAccountTokenResource,
		Scope:    nuke.Account,
		Resource: &GrafanaWorkspaceServiceAccountToken{},
		Lister:   &GrafanaWorkspaceServiceAccountTokenLister{},
	})
}

type GrafanaWorkspaceServiceAccountTokenLister struct {
	svc GrafanaV2Client
}

func (l *GrafanaWorkspaceServiceAccountTokenLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
					tokenParams := &grafana.ListWorkspaceServiceAccountTokensInput{
						WorkspaceId:      workspace.Id,
						ServiceAccountId: sa.Id,
					}
					for {
						tokenResp, err := svc.ListWorkspaceServiceAccountTokens(ctx, tokenParams)
						if err != nil {
							return nil, err
						}
						for _, token := range tokenResp.ServiceAccountTokens {
							resources = append(resources, &GrafanaWorkspaceServiceAccountToken{
								svc:              svc,
								WorkspaceID:      workspace.Id,
								ServiceAccountID: sa.Id,
								TokenID:          token.Id,
								Name:             token.Name,
							})
						}
						if tokenResp.NextToken == nil {
							break
						}
						tokenParams.NextToken = tokenResp.NextToken
					}
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

type GrafanaWorkspaceServiceAccountToken struct {
	svc              GrafanaV2Client
	WorkspaceID      *string `property:"name=WorkspaceId"`
	ServiceAccountID *string `property:"name=ServiceAccountId"`
	TokenID          *string `property:"name=TokenId"`
	Name             *string
}

func (r *GrafanaWorkspaceServiceAccountToken) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkspaceServiceAccountToken(ctx, &grafana.DeleteWorkspaceServiceAccountTokenInput{
		WorkspaceId:      r.WorkspaceID,
		ServiceAccountId: r.ServiceAccountID,
		TokenId:          r.TokenID,
	})
	return err
}

func (r *GrafanaWorkspaceServiceAccountToken) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GrafanaWorkspaceServiceAccountToken) String() string {
	return *r.Name
}
