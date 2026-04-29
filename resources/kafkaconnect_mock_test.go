package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/kafkaconnect"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockKafkaConnectClient struct {
	mock.Mock
}

func (m *mockKafkaConnectClient) ListWorkerConfigurations(ctx context.Context,
	params *kafkaconnect.ListWorkerConfigurationsInput,
	_ ...func(*kafkaconnect.Options)) (*kafkaconnect.ListWorkerConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kafkaconnect.ListWorkerConfigurationsOutput), args.Error(1)
}

func (m *mockKafkaConnectClient) DeleteWorkerConfiguration(ctx context.Context,
	params *kafkaconnect.DeleteWorkerConfigurationInput,
	_ ...func(*kafkaconnect.Options)) (*kafkaconnect.DeleteWorkerConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kafkaconnect.DeleteWorkerConfigurationOutput), args.Error(1)
}

func (m *mockKafkaConnectClient) ListCustomPlugins(ctx context.Context,
	params *kafkaconnect.ListCustomPluginsInput,
	_ ...func(*kafkaconnect.Options)) (*kafkaconnect.ListCustomPluginsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kafkaconnect.ListCustomPluginsOutput), args.Error(1)
}

func (m *mockKafkaConnectClient) DeleteCustomPlugin(ctx context.Context,
	params *kafkaconnect.DeleteCustomPluginInput,
	_ ...func(*kafkaconnect.Options)) (*kafkaconnect.DeleteCustomPluginOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*kafkaconnect.DeleteCustomPluginOutput), args.Error(1)
}

var testKafkaConnectListerOpts = &nuke.ListerOpts{}
