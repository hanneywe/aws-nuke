package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice"
)

// ApplicationdiscoveryserviceClient is the interface for the applicationdiscoveryservice SDK client methods.
type ApplicationdiscoveryserviceClient interface {
	ListConfigurations(ctx context.Context, params *applicationdiscoveryservice.ListConfigurationsInput,
		optFns ...func(*applicationdiscoveryservice.Options)) (*applicationdiscoveryservice.ListConfigurationsOutput, error)
	DeleteApplications(ctx context.Context, params *applicationdiscoveryservice.DeleteApplicationsInput,
		optFns ...func(*applicationdiscoveryservice.Options)) (*applicationdiscoveryservice.DeleteApplicationsOutput, error)
}
