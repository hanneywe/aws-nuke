package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"
)

// KafkaConnectClient is the interface for the Kafka Connect SDK client methods.
type KafkaConnectClient interface {
	ListWorkerConfigurations(ctx context.Context, params *kafkaconnect.ListWorkerConfigurationsInput,
		optFns ...func(*kafkaconnect.Options)) (*kafkaconnect.ListWorkerConfigurationsOutput, error)
	DeleteWorkerConfiguration(ctx context.Context, params *kafkaconnect.DeleteWorkerConfigurationInput,
		optFns ...func(*kafkaconnect.Options)) (*kafkaconnect.DeleteWorkerConfigurationOutput, error)
	ListCustomPlugins(ctx context.Context, params *kafkaconnect.ListCustomPluginsInput,
		optFns ...func(*kafkaconnect.Options)) (*kafkaconnect.ListCustomPluginsOutput, error)
	DeleteCustomPlugin(ctx context.Context, params *kafkaconnect.DeleteCustomPluginInput,
		optFns ...func(*kafkaconnect.Options)) (*kafkaconnect.DeleteCustomPluginOutput, error)
}
