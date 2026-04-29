package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/grafana"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GrafanaWorkspaceAPIKeyResource = "GrafanaWorkspaceApiKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     GrafanaWorkspaceAPIKeyResource,
		Scope:    nuke.Account,
		Resource: &GrafanaWorkspaceAPIKey{},
		Lister:   &GrafanaWorkspaceAPIKeyLister{},
	})
}

type GrafanaWorkspaceAPIKeyLister struct {
	svc GrafanaV2Client
}

func (l *GrafanaWorkspaceAPIKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = grafana.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	wsParams := &grafana.ListWorkspacesInput{}
	for {
		wsResp, err := svc.ListWorkspaces(ctx, wsParams)
		if err != nil {
			return nil, err
		}

		for i := range wsResp.Workspaces {
			workspaceID := wsResp.Workspaces[i].Id

			saParams := &grafana.ListWorkspaceServiceAccountsInput{
				WorkspaceId: workspaceID,
			}
			for {
				saResp, err := svc.ListWorkspaceServiceAccounts(ctx, saParams)
				if err != nil {
					break
				}

				for j := range saResp.ServiceAccounts {
					sa := &saResp.ServiceAccounts[j]

					tokenParams := &grafana.ListWorkspaceServiceAccountTokensInput{
						WorkspaceId:      workspaceID,
						ServiceAccountId: sa.Id,
					}
					for {
						tokenResp, err := svc.ListWorkspaceServiceAccountTokens(ctx, tokenParams)
						if err != nil {
							break
						}

						for k := range tokenResp.ServiceAccountTokens {
							token := &tokenResp.ServiceAccountTokens[k]
							resources = append(resources, &GrafanaWorkspaceAPIKey{
								svc:              svc,
								WorkspaceID:      workspaceID,
								ServiceAccountID: sa.Id,
								TokenID:          token.Id,
								KeyName:          token.Name,
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

		if wsResp.NextToken == nil {
			break
		}
		wsParams.NextToken = wsResp.NextToken
	}

	return resources, nil
}

type GrafanaWorkspaceAPIKey struct {
	svc              GrafanaV2Client
	WorkspaceID      *string `property:"name=WorkspaceId"`
	ServiceAccountID *string `property:"name=ServiceAccountId"`
	TokenID          *string `property:"name=TokenId"`
	KeyName          *string
}

func (r *GrafanaWorkspaceAPIKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkspaceServiceAccountToken(ctx, &grafana.DeleteWorkspaceServiceAccountTokenInput{
		WorkspaceId:      r.WorkspaceID,
		ServiceAccountId: r.ServiceAccountID,
		TokenId:          r.TokenID,
	})
	return err
}

func (r *GrafanaWorkspaceAPIKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GrafanaWorkspaceAPIKey) String() string {
	return fmt.Sprintf("%s/%s", *r.WorkspaceID, *r.KeyName)
}
