package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockElasticsearchserviceClient struct {
	mock.Mock
}

func (m *mockElasticsearchserviceClient) DescribePackages(
	ctx context.Context, params *elasticsearchservice.DescribePackagesInput,
	_ ...func(*elasticsearchservice.Options),
) (*elasticsearchservice.DescribePackagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticsearchservice.DescribePackagesOutput), args.Error(1)
}

func (m *mockElasticsearchserviceClient) DeletePackage(
	ctx context.Context, params *elasticsearchservice.DeletePackageInput,
	_ ...func(*elasticsearchservice.Options),
) (*elasticsearchservice.DeletePackageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*elasticsearchservice.DeletePackageOutput), args.Error(1)
}

var testElasticsearchserviceListerOpts = &nuke.ListerOpts{}
