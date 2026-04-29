package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/greengrass"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockGreengrassClient struct {
	mock.Mock
}

func (m *mockGreengrassClient) ListConnectorDefinitions(ctx context.Context,
	params *greengrass.ListConnectorDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListConnectorDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListConnectorDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteConnectorDefinition(ctx context.Context,
	params *greengrass.DeleteConnectorDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteConnectorDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteConnectorDefinitionOutput), args.Error(1)
}

func (m *mockGreengrassClient) ListCoreDefinitions(ctx context.Context,
	params *greengrass.ListCoreDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListCoreDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListCoreDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteCoreDefinition(ctx context.Context,
	params *greengrass.DeleteCoreDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteCoreDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteCoreDefinitionOutput), args.Error(1)
}

func (m *mockGreengrassClient) ListDeviceDefinitions(ctx context.Context,
	params *greengrass.ListDeviceDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListDeviceDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListDeviceDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteDeviceDefinition(ctx context.Context,
	params *greengrass.DeleteDeviceDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteDeviceDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteDeviceDefinitionOutput), args.Error(1)
}

func (m *mockGreengrassClient) ListFunctionDefinitions(ctx context.Context,
	params *greengrass.ListFunctionDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListFunctionDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListFunctionDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteFunctionDefinition(ctx context.Context,
	params *greengrass.DeleteFunctionDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteFunctionDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteFunctionDefinitionOutput), args.Error(1)
}

func (m *mockGreengrassClient) ListLoggerDefinitions(ctx context.Context,
	params *greengrass.ListLoggerDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListLoggerDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListLoggerDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteLoggerDefinition(ctx context.Context,
	params *greengrass.DeleteLoggerDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteLoggerDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteLoggerDefinitionOutput), args.Error(1)
}

func (m *mockGreengrassClient) ListSubscriptionDefinitions(ctx context.Context,
	params *greengrass.ListSubscriptionDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListSubscriptionDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListSubscriptionDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteSubscriptionDefinition(ctx context.Context,
	params *greengrass.DeleteSubscriptionDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteSubscriptionDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteSubscriptionDefinitionOutput), args.Error(1)
}

func (m *mockGreengrassClient) ListResourceDefinitions(ctx context.Context,
	params *greengrass.ListResourceDefinitionsInput,
	_ ...func(*greengrass.Options)) (*greengrass.ListResourceDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.ListResourceDefinitionsOutput), args.Error(1)
}

func (m *mockGreengrassClient) DeleteResourceDefinition(ctx context.Context,
	params *greengrass.DeleteResourceDefinitionInput,
	_ ...func(*greengrass.Options)) (*greengrass.DeleteResourceDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*greengrass.DeleteResourceDefinitionOutput), args.Error(1)
}

var testGreengrassListerOpts = &nuke.ListerOpts{}
