package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const KafkaConnectWorkerConfigurationResource = "KafkaConnectWorkerConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     KafkaConnectWorkerConfigurationResource,
		Scope:    nuke.Account,
		Resource: &KafkaConnectWorkerConfiguration{},
		Lister:   &KafkaConnectWorkerConfigurationLister{},
	})
}

type KafkaConnectWorkerConfigurationLister struct {
	svc KafkaConnectClient
}

func (l *KafkaConnectWorkerConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = kafkaconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := kafkaconnect.NewListWorkerConfigurationsPaginator(svc, &kafkaconnect.ListWorkerConfigurationsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, wc := range resp.WorkerConfigurations {
			resources = append(resources, &KafkaConnectWorkerConfiguration{
				svc:                    svc,
				WorkerConfigurationArn: wc.WorkerConfigurationArn,
				Name:                   wc.Name,
			})
		}
	}

	return resources, nil
}

type KafkaConnectWorkerConfiguration struct {
	svc                    KafkaConnectClient
	WorkerConfigurationArn *string
	Name                   *string
}

func (r *KafkaConnectWorkerConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkerConfiguration(ctx, &kafkaconnect.DeleteWorkerConfigurationInput{
		WorkerConfigurationArn: r.WorkerConfigurationArn,
	})
	return err
}

func (r *KafkaConnectWorkerConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *KafkaConnectWorkerConfiguration) String() string {
	return *r.Name
}
