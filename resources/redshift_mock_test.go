package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
)

type mockRedshiftClient struct {
	mock.Mock
}

func (m *mockRedshiftClient) DescribeHsmClientCertificates(
	ctx context.Context, params *redshift.DescribeHsmClientCertificatesInput,
	_ ...func(*redshift.Options),
) (*redshift.DescribeHsmClientCertificatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DescribeHsmClientCertificatesOutput), args.Error(1)
}

func (m *mockRedshiftClient) DeleteHsmClientCertificate(
	ctx context.Context, params *redshift.DeleteHsmClientCertificateInput,
	_ ...func(*redshift.Options),
) (*redshift.DeleteHsmClientCertificateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DeleteHsmClientCertificateOutput), args.Error(1)
}

func (m *mockRedshiftClient) DescribeSnapshotCopyGrants(
	ctx context.Context, params *redshift.DescribeSnapshotCopyGrantsInput,
	_ ...func(*redshift.Options),
) (*redshift.DescribeSnapshotCopyGrantsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DescribeSnapshotCopyGrantsOutput), args.Error(1)
}

func (m *mockRedshiftClient) DeleteSnapshotCopyGrant(
	ctx context.Context, params *redshift.DeleteSnapshotCopyGrantInput,
	_ ...func(*redshift.Options),
) (*redshift.DeleteSnapshotCopyGrantOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DeleteSnapshotCopyGrantOutput), args.Error(1)
}

func (m *mockRedshiftClient) DescribeEventSubscriptions(
	ctx context.Context, params *redshift.DescribeEventSubscriptionsInput,
	_ ...func(*redshift.Options),
) (*redshift.DescribeEventSubscriptionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DescribeEventSubscriptionsOutput), args.Error(1)
}

func (m *mockRedshiftClient) DeleteEventSubscription(
	ctx context.Context, params *redshift.DeleteEventSubscriptionInput,
	_ ...func(*redshift.Options),
) (*redshift.DeleteEventSubscriptionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DeleteEventSubscriptionOutput), args.Error(1)
}

func (m *mockRedshiftClient) DescribeHsmConfigurations(
	ctx context.Context, params *redshift.DescribeHsmConfigurationsInput,
	_ ...func(*redshift.Options),
) (*redshift.DescribeHsmConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DescribeHsmConfigurationsOutput), args.Error(1)
}

func (m *mockRedshiftClient) DeleteHsmConfiguration(
	ctx context.Context, params *redshift.DeleteHsmConfigurationInput,
	_ ...func(*redshift.Options),
) (*redshift.DeleteHsmConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*redshift.DeleteHsmConfigurationOutput), args.Error(1)
}
