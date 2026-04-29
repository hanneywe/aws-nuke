package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivschat"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSChatLoggingConfigurationResource = "IVSChatLoggingConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSChatLoggingConfigurationResource,
		Scope:    nuke.Account,
		Resource: &IVSChatLoggingConfiguration{},
		Lister:   &IVSChatLoggingConfigurationLister{},
	})
}

type IVSChatLoggingConfigurationLister struct {
	svc IVSChatClient
}

func (l *IVSChatLoggingConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ivschat.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ivschat.NewListLoggingConfigurationsPaginator(svc, &ivschat.ListLoggingConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, loggingConfig := range resp.LoggingConfigurations {
			resources = append(resources, &IVSChatLoggingConfiguration{
				svc:  svc,
				Arn:  loggingConfig.Arn,
				Name: loggingConfig.Name,
				Tags: loggingConfig.Tags,
			})
		}
	}

	return resources, nil
}

type IVSChatLoggingConfiguration struct {
	svc  IVSChatClient
	Arn  *string
	Name *string
	Tags map[string]string
}

func (r *IVSChatLoggingConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLoggingConfiguration(ctx, &ivschat.DeleteLoggingConfigurationInput{
		Identifier: r.Arn,
	})
	return err
}

func (r *IVSChatLoggingConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSChatLoggingConfiguration) String() string {
	return *r.Name
}
