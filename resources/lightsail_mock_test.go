package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockLightsailClient struct {
	mock.Mock
}

func (m *mockLightsailClient) GetBuckets(ctx context.Context,
	params *lightsail.GetBucketsInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetBucketsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetBucketsOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteBucket(ctx context.Context,
	params *lightsail.DeleteBucketInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteBucketOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteBucketOutput), args.Error(1)
}

func (m *mockLightsailClient) GetCertificates(ctx context.Context,
	params *lightsail.GetCertificatesInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetCertificatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetCertificatesOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteCertificate(ctx context.Context,
	params *lightsail.DeleteCertificateInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteCertificateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteCertificateOutput), args.Error(1)
}

func (m *mockLightsailClient) GetContainerServices(ctx context.Context,
	params *lightsail.GetContainerServicesInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetContainerServicesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetContainerServicesOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteContainerService(ctx context.Context,
	params *lightsail.DeleteContainerServiceInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteContainerServiceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteContainerServiceOutput), args.Error(1)
}

func (m *mockLightsailClient) GetContactMethods(ctx context.Context,
	params *lightsail.GetContactMethodsInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetContactMethodsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetContactMethodsOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteContactMethod(ctx context.Context,
	params *lightsail.DeleteContactMethodInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteContactMethodOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteContactMethodOutput), args.Error(1)
}

func (m *mockLightsailClient) GetBucketAccessKeys(ctx context.Context,
	params *lightsail.GetBucketAccessKeysInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetBucketAccessKeysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetBucketAccessKeysOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteBucketAccessKey(ctx context.Context,
	params *lightsail.DeleteBucketAccessKeyInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteBucketAccessKeyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteBucketAccessKeyOutput), args.Error(1)
}

func (m *mockLightsailClient) GetDiskSnapshots(ctx context.Context,
	params *lightsail.GetDiskSnapshotsInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetDiskSnapshotsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetDiskSnapshotsOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteDiskSnapshot(ctx context.Context,
	params *lightsail.DeleteDiskSnapshotInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteDiskSnapshotOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteDiskSnapshotOutput), args.Error(1)
}

func (m *mockLightsailClient) GetRelationalDatabases(ctx context.Context,
	params *lightsail.GetRelationalDatabasesInput,
	_ ...func(*lightsail.Options)) (*lightsail.GetRelationalDatabasesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.GetRelationalDatabasesOutput), args.Error(1)
}

func (m *mockLightsailClient) DeleteRelationalDatabase(ctx context.Context,
	params *lightsail.DeleteRelationalDatabaseInput,
	_ ...func(*lightsail.Options)) (*lightsail.DeleteRelationalDatabaseOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*lightsail.DeleteRelationalDatabaseOutput), args.Error(1)
}

var testLightsailListerOpts = &nuke.ListerOpts{}
