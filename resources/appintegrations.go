package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appintegrations"
)

// AppIntegrationsClient is an interface for the AppIntegrations SDK client methods used by all AppIntegrations resources.
// It enables mock testing of List and Remove operations.
type AppIntegrationsClient interface {
	ListEventIntegrations(ctx context.Context, params *appintegrations.ListEventIntegrationsInput,
		optFns ...func(*appintegrations.Options)) (*appintegrations.ListEventIntegrationsOutput, error)
	DeleteEventIntegration(ctx context.Context, params *appintegrations.DeleteEventIntegrationInput,
		optFns ...func(*appintegrations.Options)) (*appintegrations.DeleteEventIntegrationOutput, error)
}
