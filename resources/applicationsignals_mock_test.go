package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"
)

type mockApplicationSignalsClient struct {
	mock.Mock
}

func (m *mockApplicationSignalsClient) ListGroupingAttributeDefinitions(
	ctx context.Context, params *applicationsignals.ListGroupingAttributeDefinitionsInput,
	_ ...func(*applicationsignals.Options),
) (*applicationsignals.ListGroupingAttributeDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*applicationsignals.ListGroupingAttributeDefinitionsOutput), args.Error(1)
}

func (m *mockApplicationSignalsClient) DeleteGroupingConfiguration(
	ctx context.Context, params *applicationsignals.DeleteGroupingConfigurationInput,
	_ ...func(*applicationsignals.Options),
) (*applicationsignals.DeleteGroupingConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*applicationsignals.DeleteGroupingConfigurationOutput), args.Error(1)
}
