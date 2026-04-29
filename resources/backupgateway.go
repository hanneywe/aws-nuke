package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/backupgateway"
)

// BackupGatewayClient is an interface for the Backup Gateway SDK client methods used by all Backup Gateway resources.
// It enables mock testing of List and Remove operations.
type BackupGatewayClient interface {
	ListHypervisors(ctx context.Context, params *backupgateway.ListHypervisorsInput,
		optFns ...func(*backupgateway.Options)) (*backupgateway.ListHypervisorsOutput, error)
	DeleteHypervisor(ctx context.Context, params *backupgateway.DeleteHypervisorInput,
		optFns ...func(*backupgateway.Options)) (*backupgateway.DeleteHypervisorOutput, error)
}
