package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTWirelessNetworkAnalyzerConfigurationResource = "IoTWirelessNetworkAnalyzerConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTWirelessNetworkAnalyzerConfigurationResource,
		Scope:    nuke.Account,
		Resource: &IoTWirelessNetworkAnalyzerConfiguration{},
		Lister:   &IoTWirelessNetworkAnalyzerConfigurationLister{},
	})
}

type IoTWirelessNetworkAnalyzerConfigurationLister struct {
	svc IoTWirelessClient
}

func (l *IoTWirelessNetworkAnalyzerConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iotwireless.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iotwireless.NewListNetworkAnalyzerConfigurationsPaginator(svc, &iotwireless.ListNetworkAnalyzerConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, cfg := range resp.NetworkAnalyzerConfigurationList {
			resources = append(resources, &IoTWirelessNetworkAnalyzerConfiguration{
				svc:  svc,
				Name: cfg.Name,
			})
		}
	}

	return resources, nil
}

type IoTWirelessNetworkAnalyzerConfiguration struct {
	svc  IoTWirelessClient
	Name *string
}

func (r *IoTWirelessNetworkAnalyzerConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteNetworkAnalyzerConfiguration(ctx, &iotwireless.DeleteNetworkAnalyzerConfigurationInput{
		ConfigurationName: r.Name,
	})
	return err
}

func (r *IoTWirelessNetworkAnalyzerConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTWirelessNetworkAnalyzerConfiguration) String() string {
	return *r.Name
}
