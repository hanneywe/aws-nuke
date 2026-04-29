package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/inspector2"
)

// Inspector2V2Client is the interface for the Inspector2 SDK v2 client methods used by new resources.
type Inspector2V2Client interface {
	ListCodeSecurityScanConfigurations(ctx context.Context, params *inspector2.ListCodeSecurityScanConfigurationsInput,
		optFns ...func(*inspector2.Options)) (*inspector2.ListCodeSecurityScanConfigurationsOutput, error)
	DeleteCodeSecurityScanConfiguration(ctx context.Context, params *inspector2.DeleteCodeSecurityScanConfigurationInput,
		optFns ...func(*inspector2.Options)) (*inspector2.DeleteCodeSecurityScanConfigurationOutput, error)
}
