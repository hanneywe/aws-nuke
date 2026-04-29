package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const KafkaConnectCustomPluginResource = "KafkaConnectCustomPlugin"

func init() {
	registry.Register(&registry.Registration{
		Name:     KafkaConnectCustomPluginResource,
		Scope:    nuke.Account,
		Resource: &KafkaConnectCustomPlugin{},
		Lister:   &KafkaConnectCustomPluginLister{},
	})
}

type KafkaConnectCustomPluginLister struct {
	svc KafkaConnectClient
}

func (l *KafkaConnectCustomPluginLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = kafkaconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &kafkaconnect.ListCustomPluginsInput{}
	for {
		resp, err := svc.ListCustomPlugins(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, cp := range resp.CustomPlugins {
			resources = append(resources, &KafkaConnectCustomPlugin{
				svc:             svc,
				CustomPluginArn: cp.CustomPluginArn,
				Name:            cp.Name,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type KafkaConnectCustomPlugin struct {
	svc             KafkaConnectClient
	CustomPluginArn *string
	Name            *string
}

func (r *KafkaConnectCustomPlugin) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCustomPlugin(ctx, &kafkaconnect.DeleteCustomPluginInput{
		CustomPluginArn: r.CustomPluginArn,
	})
	return err
}

func (r *KafkaConnectCustomPlugin) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *KafkaConnectCustomPlugin) String() string {
	return *r.Name
}
