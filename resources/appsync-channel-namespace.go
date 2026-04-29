package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appsync"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppSyncChannelNamespaceResource = "AppSyncChannelNamespace"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppSyncChannelNamespaceResource,
		Scope:    nuke.Account,
		Resource: &AppSyncChannelNamespace{},
		Lister:   &AppSyncChannelNamespaceLister{},
	})
}

type AppSyncChannelNamespaceLister struct {
	svc AppSyncClient
}

func (l *AppSyncChannelNamespaceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = appsync.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	apiPaginator := appsync.NewListApisPaginator(svc, &appsync.ListApisInput{})
	for apiPaginator.HasMorePages() {
		apiResp, err := apiPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, api := range apiResp.Apis {
			nsPaginator := appsync.NewListChannelNamespacesPaginator(svc, &appsync.ListChannelNamespacesInput{
				ApiId: api.ApiId,
			})
			for nsPaginator.HasMorePages() {
				nsResp, err := nsPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, ns := range nsResp.ChannelNamespaces {
					resources = append(resources, &AppSyncChannelNamespace{
						svc:   svc,
						APIID: api.ApiId,
						Name:  ns.Name,
						Tags:  ns.Tags,
					})
				}
			}
		}
	}

	return resources, nil
}

type AppSyncChannelNamespace struct {
	svc   AppSyncClient
	APIID *string `property:"name=ApiId"`
	Name  *string
	Tags  map[string]string
}

func (r *AppSyncChannelNamespace) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
		ApiId: r.APIID,
		Name:  r.Name,
	})
	return err
}

func (r *AppSyncChannelNamespace) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppSyncChannelNamespace) String() string {
	return *r.Name
}
