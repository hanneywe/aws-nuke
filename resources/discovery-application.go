package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice"
	applicationdiscoveryservicetypes "github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const DiscoveryApplicationResource = "DiscoveryApplication"

func init() {
	registry.Register(&registry.Registration{
		Name:     DiscoveryApplicationResource,
		Scope:    nuke.Account,
		Resource: &DiscoveryApplication{},
		Lister:   &DiscoveryApplicationLister{},
	})
}

type DiscoveryApplicationLister struct {
	svc ApplicationdiscoveryserviceClient
}

func (l *DiscoveryApplicationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = applicationdiscoveryservice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := applicationdiscoveryservice.NewListConfigurationsPaginator(svc, &applicationdiscoveryservice.ListConfigurationsInput{
		ConfigurationType: applicationdiscoveryservicetypes.ConfigurationItemTypeApplication,
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Configurations {
			resources = append(resources, &DiscoveryApplication{
				svc:             svc,
				ConfigurationID: aws.String(item["application.configurationId"]),
				Name:            aws.String(item["application.name"]),
			})
		}
	}

	return resources, nil
}

type DiscoveryApplication struct {
	svc             ApplicationdiscoveryserviceClient
	ConfigurationID *string `property:"name=ConfigurationId"`
	Name            *string
}

func (r *DiscoveryApplication) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteApplications(ctx, &applicationdiscoveryservice.DeleteApplicationsInput{
		ConfigurationIds: []string{aws.ToString(r.ConfigurationID)},
	})
	return err
}

func (r *DiscoveryApplication) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *DiscoveryApplication) String() string {
	return *r.ConfigurationID
}
