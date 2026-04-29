package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appsync"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppSyncAPIKeyResource = "AppSyncAPIKey" //nolint:gosec // G101 - not a credential, this is a resource type name

func init() {
	registry.Register(&registry.Registration{
		Name:     AppSyncAPIKeyResource,
		Scope:    nuke.Account,
		Resource: &AppSyncAPIKey{},
		Lister:   &AppSyncAPIKeyLister{},
	})
}

type AppSyncAPIKeyLister struct {
	svc AppSyncClient
}

func (l *AppSyncAPIKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = appsync.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	apiParams := &appsync.ListGraphqlApisInput{}
	for {
		apiResp, err := svc.ListGraphqlApis(ctx, apiParams)
		if err != nil {
			return nil, err
		}
		for i := range apiResp.GraphqlApis {
			keyParams := &appsync.ListApiKeysInput{ApiId: apiResp.GraphqlApis[i].ApiId}
			for {
				keyResp, err := svc.ListApiKeys(ctx, keyParams)
				if err != nil {
					return nil, err
				}
				for _, key := range keyResp.ApiKeys {
					resources = append(resources, &AppSyncAPIKey{
						svc:   svc,
						ID:    key.Id,
						APIID: apiResp.GraphqlApis[i].ApiId,
					})
				}
				if keyResp.NextToken == nil {
					break
				}
				keyParams.NextToken = keyResp.NextToken
			}
		}
		if apiResp.NextToken == nil {
			break
		}
		apiParams.NextToken = apiResp.NextToken
	}

	return resources, nil
}

type AppSyncAPIKey struct {
	svc   AppSyncClient
	ID    *string `property:"name=Id"`
	APIID *string `property:"name=ApiId"`
}

func (r *AppSyncAPIKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
		ApiId: r.APIID,
		Id:    r.ID,
	})
	return err
}

func (r *AppSyncAPIKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppSyncAPIKey) String() string {
	return *r.ID
}
