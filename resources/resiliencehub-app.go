package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resiliencehub"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ResilienceHubAppResource = "ResilienceHubApp"

func init() {
	registry.Register(&registry.Registration{
		Name:     ResilienceHubAppResource,
		Scope:    nuke.Account,
		Resource: &ResilienceHubApp{},
		Lister:   &ResilienceHubAppLister{},
	})
}

type ResilienceHubAppLister struct {
	svc ResilienceHubClient
}

func (l *ResilienceHubAppLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = resiliencehub.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := resiliencehub.NewListAppsPaginator(svc, &resiliencehub.ListAppsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range resp.AppSummaries {
			resources = append(resources, &ResilienceHubApp{
				svc:    svc,
				AppArn: resp.AppSummaries[i].AppArn,
				Name:   resp.AppSummaries[i].Name,
			})
		}
	}
	return resources, nil
}

type ResilienceHubApp struct {
	svc    ResilienceHubClient
	AppArn *string
	Name   *string
}

func (r *ResilienceHubApp) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteApp(ctx, &resiliencehub.DeleteAppInput{
		AppArn: r.AppArn,
	})
	return err
}

func (r *ResilienceHubApp) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ResilienceHubApp) String() string {
	return *r.Name
}
