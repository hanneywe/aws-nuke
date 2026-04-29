package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotsitewise"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockIotsitewiseClient struct {
	mock.Mock
}

func (m *mockIotsitewiseClient) ListAssetModels(
	ctx context.Context, params *iotsitewise.ListAssetModelsInput,
	_ ...func(*iotsitewise.Options),
) (*iotsitewise.ListAssetModelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotsitewise.ListAssetModelsOutput), args.Error(1)
}

func (m *mockIotsitewiseClient) ListAssetModelCompositeModels(
	ctx context.Context, params *iotsitewise.ListAssetModelCompositeModelsInput,
	_ ...func(*iotsitewise.Options),
) (*iotsitewise.ListAssetModelCompositeModelsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotsitewise.ListAssetModelCompositeModelsOutput), args.Error(1)
}

func (m *mockIotsitewiseClient) DeleteAssetModelCompositeModel(
	ctx context.Context, params *iotsitewise.DeleteAssetModelCompositeModelInput,
	_ ...func(*iotsitewise.Options),
) (*iotsitewise.DeleteAssetModelCompositeModelOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotsitewise.DeleteAssetModelCompositeModelOutput), args.Error(1)
}

func (m *mockIotsitewiseClient) ListDatasets(
	ctx context.Context, params *iotsitewise.ListDatasetsInput,
	_ ...func(*iotsitewise.Options),
) (*iotsitewise.ListDatasetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotsitewise.ListDatasetsOutput), args.Error(1)
}

func (m *mockIotsitewiseClient) DeleteDataset(
	ctx context.Context, params *iotsitewise.DeleteDatasetInput,
	_ ...func(*iotsitewise.Options),
) (*iotsitewise.DeleteDatasetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*iotsitewise.DeleteDatasetOutput), args.Error(1)
}

var testIotsitewiseListerOpts = &nuke.ListerOpts{}
