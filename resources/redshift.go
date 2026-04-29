package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

// RedshiftClient is the interface for the Redshift SDK v2 client methods
// used by new Redshift sub-resources. Existing Redshift resources use SDK v1.
type RedshiftClient interface {
	DescribeHsmClientCertificates(ctx context.Context, params *redshift.DescribeHsmClientCertificatesInput,
		optFns ...func(*redshift.Options)) (*redshift.DescribeHsmClientCertificatesOutput, error)
	DeleteHsmClientCertificate(ctx context.Context, params *redshift.DeleteHsmClientCertificateInput,
		optFns ...func(*redshift.Options)) (*redshift.DeleteHsmClientCertificateOutput, error)
	DescribeSnapshotCopyGrants(ctx context.Context, params *redshift.DescribeSnapshotCopyGrantsInput,
		optFns ...func(*redshift.Options)) (*redshift.DescribeSnapshotCopyGrantsOutput, error)
	DeleteSnapshotCopyGrant(ctx context.Context, params *redshift.DeleteSnapshotCopyGrantInput,
		optFns ...func(*redshift.Options)) (*redshift.DeleteSnapshotCopyGrantOutput, error)
	DescribeEventSubscriptions(ctx context.Context, params *redshift.DescribeEventSubscriptionsInput,
		optFns ...func(*redshift.Options)) (*redshift.DescribeEventSubscriptionsOutput, error)
	DeleteEventSubscription(ctx context.Context, params *redshift.DeleteEventSubscriptionInput,
		optFns ...func(*redshift.Options)) (*redshift.DeleteEventSubscriptionOutput, error)
	DescribeHsmConfigurations(ctx context.Context, params *redshift.DescribeHsmConfigurationsInput,
		optFns ...func(*redshift.Options)) (*redshift.DescribeHsmConfigurationsOutput, error)
	DeleteHsmConfiguration(ctx context.Context, params *redshift.DeleteHsmConfigurationInput,
		optFns ...func(*redshift.Options)) (*redshift.DeleteHsmConfigurationOutput, error)
}
