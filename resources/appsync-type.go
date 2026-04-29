package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appsync"
	appsynctypes "github.com/aws/aws-sdk-go-v2/service/appsync/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppSyncTypeResource = "AppSyncType"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppSyncTypeResource,
		Scope:    nuke.Account,
		Resource: &AppSyncType{},
		Lister:   &AppSyncTypeLister{},
	})
}

type AppSyncTypeLister struct {
	svc AppSyncClient
}

func (l *AppSyncTypeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			typeParams := &appsync.ListTypesInput{
				ApiId:  apiResp.GraphqlApis[i].ApiId,
				Format: appsynctypes.TypeDefinitionFormatSdl,
			}
			for {
				typeResp, err := svc.ListTypes(ctx, typeParams)
				if err != nil {
					return nil, err
				}
				for _, t := range typeResp.Types {
					resources = append(resources, &AppSyncType{
						svc:      svc,
						TypeName: t.Name,
						APIID:    apiResp.GraphqlApis[i].ApiId,
					})
				}
				if typeResp.NextToken == nil {
					break
				}
				typeParams.NextToken = typeResp.NextToken
			}
		}
		if apiResp.NextToken == nil {
			break
		}
		apiParams.NextToken = apiResp.NextToken
	}

	return resources, nil
}

type AppSyncType struct {
	svc      AppSyncClient
	TypeName *string
	APIID    *string `property:"name=ApiId"`
}

func (r *AppSyncType) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteType(ctx, &appsync.DeleteTypeInput{
		ApiId:    r.APIID,
		TypeName: r.TypeName,
	})
	return err
}

func (r *AppSyncType) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppSyncType) String() string {
	return *r.TypeName
}
