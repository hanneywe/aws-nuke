package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/docdb"
)

type mockDocDBV2Client struct {
	mock.Mock
}

func (m *mockDocDBV2Client) DescribeGlobalClusters(ctx context.Context, params *docdb.DescribeGlobalClustersInput,
	_ ...func(*docdb.Options)) (*docdb.DescribeGlobalClustersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*docdb.DescribeGlobalClustersOutput), args.Error(1)
}

func (m *mockDocDBV2Client) DeleteGlobalCluster(ctx context.Context, params *docdb.DeleteGlobalClusterInput,
	_ ...func(*docdb.Options)) (*docdb.DeleteGlobalClusterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*docdb.DeleteGlobalClusterOutput), args.Error(1)
}

func (m *mockDocDBV2Client) ModifyGlobalCluster(ctx context.Context, params *docdb.ModifyGlobalClusterInput,
	_ ...func(*docdb.Options)) (*docdb.ModifyGlobalClusterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*docdb.ModifyGlobalClusterOutput), args.Error(1)
}
