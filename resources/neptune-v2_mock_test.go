package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/neptune"
)

type mockNeptuneV2Client struct {
	mock.Mock
}

func (m *mockNeptuneV2Client) DescribeDBParameterGroups(ctx context.Context, params *neptune.DescribeDBParameterGroupsInput,
	_ ...func(*neptune.Options)) (*neptune.DescribeDBParameterGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DescribeDBParameterGroupsOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DeleteDBParameterGroup(ctx context.Context, params *neptune.DeleteDBParameterGroupInput,
	_ ...func(*neptune.Options)) (*neptune.DeleteDBParameterGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DeleteDBParameterGroupOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DescribeDBClusterParameterGroups(ctx context.Context, params *neptune.DescribeDBClusterParameterGroupsInput,
	_ ...func(*neptune.Options)) (*neptune.DescribeDBClusterParameterGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DescribeDBClusterParameterGroupsOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DeleteDBClusterParameterGroup(ctx context.Context, params *neptune.DeleteDBClusterParameterGroupInput,
	_ ...func(*neptune.Options)) (*neptune.DeleteDBClusterParameterGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DeleteDBClusterParameterGroupOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DescribeEventSubscriptions(ctx context.Context, params *neptune.DescribeEventSubscriptionsInput,
	_ ...func(*neptune.Options)) (*neptune.DescribeEventSubscriptionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DescribeEventSubscriptionsOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DeleteEventSubscription(ctx context.Context, params *neptune.DeleteEventSubscriptionInput,
	_ ...func(*neptune.Options)) (*neptune.DeleteEventSubscriptionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DeleteEventSubscriptionOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DescribeGlobalClusters(ctx context.Context, params *neptune.DescribeGlobalClustersInput,
	_ ...func(*neptune.Options)) (*neptune.DescribeGlobalClustersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DescribeGlobalClustersOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DeleteGlobalCluster(ctx context.Context, params *neptune.DeleteGlobalClusterInput,
	_ ...func(*neptune.Options)) (*neptune.DeleteGlobalClusterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DeleteGlobalClusterOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) ModifyGlobalCluster(ctx context.Context, params *neptune.ModifyGlobalClusterInput,
	_ ...func(*neptune.Options)) (*neptune.ModifyGlobalClusterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.ModifyGlobalClusterOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DescribeDBClusterEndpoints(ctx context.Context, params *neptune.DescribeDBClusterEndpointsInput,
	_ ...func(*neptune.Options)) (*neptune.DescribeDBClusterEndpointsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DescribeDBClusterEndpointsOutput), args.Error(1)
}

func (m *mockNeptuneV2Client) DeleteDBClusterEndpoint(ctx context.Context, params *neptune.DeleteDBClusterEndpointInput,
	_ ...func(*neptune.Options)) (*neptune.DeleteDBClusterEndpointOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*neptune.DeleteDBClusterEndpointOutput), args.Error(1)
}
