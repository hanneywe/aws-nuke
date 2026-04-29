package resources

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/smithy-go"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppSyncAPICacheResource = "AppSyncAPICache"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppSyncAPICacheResource,
		Scope:    nuke.Account,
		Resource: &AppSyncAPICache{},
		Lister:   &AppSyncAPICacheLister{},
	})
}

type AppSyncAPICacheLister struct {
	svc AppSyncClient
}

func (l *AppSyncAPICacheLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = appsync.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// Iterate over all GraphQL APIs
	apiParams := &appsync.ListGraphqlApisInput{}
	for {
		apiResp, err := svc.ListGraphqlApis(ctx, apiParams)
		if err != nil {
			return nil, err
		}
		for i := range apiResp.GraphqlApis {
			_, err := svc.GetApiCache(ctx, &appsync.GetApiCacheInput{
				ApiId: apiResp.GraphqlApis[i].ApiId,
			})
			if err != nil {
				var apiErr smithy.APIError
				if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFoundException" {
					continue
				}
				return nil, err
			}
			resources = append(resources, &AppSyncAPICache{
				svc:   svc,
				APIID: apiResp.GraphqlApis[i].ApiId,
			})
		}
		if apiResp.NextToken == nil {
			break
		}
		apiParams.NextToken = apiResp.NextToken
	}

	return resources, nil
}

type AppSyncAPICache struct {
	svc   AppSyncClient
	APIID *string `property:"name=ApiId"`
}

func (r *AppSyncAPICache) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteApiCache(ctx, &appsync.DeleteApiCacheInput{
		ApiId: r.APIID,
	})
	return err
}

func (r *AppSyncAPICache) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppSyncAPICache) String() string {
	return *r.APIID
}
