package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/efs"
)

type mockEFSV2Client struct {
	mock.Mock
}

func (m *mockEFSV2Client) DescribeAccessPoints(ctx context.Context, params *efs.DescribeAccessPointsInput,
	_ ...func(*efs.Options)) (*efs.DescribeAccessPointsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DescribeAccessPointsOutput), args.Error(1)
}

func (m *mockEFSV2Client) DeleteAccessPoint(ctx context.Context, params *efs.DeleteAccessPointInput,
	_ ...func(*efs.Options)) (*efs.DeleteAccessPointOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DeleteAccessPointOutput), args.Error(1)
}

func (m *mockEFSV2Client) DescribeFileSystems(ctx context.Context, params *efs.DescribeFileSystemsInput,
	_ ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DescribeFileSystemsOutput), args.Error(1)
}

func (m *mockEFSV2Client) DescribeBackupPolicy(ctx context.Context, params *efs.DescribeBackupPolicyInput,
	_ ...func(*efs.Options)) (*efs.DescribeBackupPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DescribeBackupPolicyOutput), args.Error(1)
}

func (m *mockEFSV2Client) PutBackupPolicy(ctx context.Context, params *efs.PutBackupPolicyInput,
	_ ...func(*efs.Options)) (*efs.PutBackupPolicyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.PutBackupPolicyOutput), args.Error(1)
}

func (m *mockEFSV2Client) ListTagsForResource(ctx context.Context, params *efs.ListTagsForResourceInput,
	_ ...func(*efs.Options)) (*efs.ListTagsForResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.ListTagsForResourceOutput), args.Error(1)
}

func (m *mockEFSV2Client) UntagResource(ctx context.Context, params *efs.UntagResourceInput,
	_ ...func(*efs.Options)) (*efs.UntagResourceOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.UntagResourceOutput), args.Error(1)
}

func (m *mockEFSV2Client) DescribeLifecycleConfiguration(ctx context.Context, params *efs.DescribeLifecycleConfigurationInput,
	_ ...func(*efs.Options)) (*efs.DescribeLifecycleConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DescribeLifecycleConfigurationOutput), args.Error(1)
}

func (m *mockEFSV2Client) PutLifecycleConfiguration(ctx context.Context, params *efs.PutLifecycleConfigurationInput,
	_ ...func(*efs.Options)) (*efs.PutLifecycleConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.PutLifecycleConfigurationOutput), args.Error(1)
}

func (m *mockEFSV2Client) DescribeReplicationConfigurations(ctx context.Context, params *efs.DescribeReplicationConfigurationsInput,
	_ ...func(*efs.Options)) (*efs.DescribeReplicationConfigurationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DescribeReplicationConfigurationsOutput), args.Error(1)
}

func (m *mockEFSV2Client) DeleteReplicationConfiguration(ctx context.Context, params *efs.DeleteReplicationConfigurationInput,
	_ ...func(*efs.Options)) (*efs.DeleteReplicationConfigurationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*efs.DeleteReplicationConfigurationOutput), args.Error(1)
}
