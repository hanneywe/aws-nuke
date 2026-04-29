package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/databasemigrationservice"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockDMSClient struct {
	mock.Mock
}

func (m *mockDMSClient) DescribeDataProviders(
	ctx context.Context, params *databasemigrationservice.DescribeDataProvidersInput,
	_ ...func(*databasemigrationservice.Options),
) (*databasemigrationservice.DescribeDataProvidersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*databasemigrationservice.DescribeDataProvidersOutput), args.Error(1)
}

func (m *mockDMSClient) DeleteDataProvider(
	ctx context.Context, params *databasemigrationservice.DeleteDataProviderInput,
	_ ...func(*databasemigrationservice.Options),
) (*databasemigrationservice.DeleteDataProviderOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*databasemigrationservice.DeleteDataProviderOutput), args.Error(1)
}

var testDMSListerOpts = &nuke.ListerOpts{}
