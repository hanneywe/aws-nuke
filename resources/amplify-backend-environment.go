package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/amplify"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AmplifyBackendEnvironmentResource = "AmplifyBackendEnvironment"

func init() {
	registry.Register(&registry.Registration{
		Name:     AmplifyBackendEnvironmentResource,
		Scope:    nuke.Account,
		Resource: &AmplifyBackendEnvironment{},
		Lister:   &AmplifyBackendEnvironmentLister{},
	})
}

type AmplifyBackendEnvironmentLister struct {
	svc AmplifyClient
}

func (l *AmplifyBackendEnvironmentLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = amplify.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	appPaginator := amplify.NewListAppsPaginator(svc, &amplify.ListAppsInput{})
	for appPaginator.HasMorePages() {
		appResp, err := appPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range appResp.Apps {
			params := &amplify.ListBackendEnvironmentsInput{AppId: appResp.Apps[i].AppId}
			for {
				resp, err := svc.ListBackendEnvironments(ctx, params)
				if err != nil {
					return nil, err
				}
				for _, env := range resp.BackendEnvironments {
					resources = append(resources, &AmplifyBackendEnvironment{
						svc:             svc,
						EnvironmentName: env.EnvironmentName,
						AppID:           appResp.Apps[i].AppId,
					})
				}
				if resp.NextToken == nil {
					break
				}
				params.NextToken = resp.NextToken
			}
		}
	}
	return resources, nil
}

type AmplifyBackendEnvironment struct {
	svc             AmplifyClient
	EnvironmentName *string
	AppID           *string `property:"name=AppId"`
}

func (r *AmplifyBackendEnvironment) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBackendEnvironment(ctx, &amplify.DeleteBackendEnvironmentInput{
		AppId:           r.AppID,
		EnvironmentName: r.EnvironmentName,
	})
	return err
}

func (r *AmplifyBackendEnvironment) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AmplifyBackendEnvironment) String() string {
	return *r.EnvironmentName
}
