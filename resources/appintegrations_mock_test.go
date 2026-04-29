package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/appintegrations"
)

type mockAppIntegrationsClient struct {
	mock.Mock
}

func (m *mockAppIntegrationsClient) ListEventIntegrations(ctx context.Context, params *appintegrations.ListEventIntegrationsInput,
	_ ...func(*appintegrations.Options)) (*appintegrations.ListEventIntegrationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appintegrations.ListEventIntegrationsOutput), args.Error(1)
}

func (m *mockAppIntegrationsClient) DeleteEventIntegration(ctx context.Context, params *appintegrations.DeleteEventIntegrationInput,
	_ ...func(*appintegrations.Options)) (*appintegrations.DeleteEventIntegrationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*appintegrations.DeleteEventIntegrationOutput), args.Error(1)
}
