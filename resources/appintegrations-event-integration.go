package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appintegrations"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppIntegrationsEventIntegrationResource = "AppIntegrationsEventIntegration"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppIntegrationsEventIntegrationResource,
		Scope:    nuke.Account,
		Resource: &AppIntegrationsEventIntegration{},
		Lister:   &AppIntegrationsEventIntegrationLister{},
	})
}

type AppIntegrationsEventIntegrationLister struct {
	svc AppIntegrationsClient
}

func (l *AppIntegrationsEventIntegrationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = appintegrations.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := appintegrations.NewListEventIntegrationsPaginator(svc, &appintegrations.ListEventIntegrationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.EventIntegrations {
			resources = append(resources, &AppIntegrationsEventIntegration{
				svc:                 svc,
				Name:                item.Name,
				EventIntegrationArn: item.EventIntegrationArn,
			})
		}
	}

	return resources, nil
}

type AppIntegrationsEventIntegration struct {
	svc                 AppIntegrationsClient
	Name                *string
	EventIntegrationArn *string
}

func (r *AppIntegrationsEventIntegration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventIntegration(ctx, &appintegrations.DeleteEventIntegrationInput{
		Name: r.Name,
	})
	return err
}

func (r *AppIntegrationsEventIntegration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppIntegrationsEventIntegration) String() string {
	return *r.Name
}
