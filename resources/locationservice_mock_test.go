package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/location"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockLocationServiceClient struct {
	mock.Mock
}

func (m *mockLocationServiceClient) ListMaps(ctx context.Context,
	params *location.ListMapsInput,
	_ ...func(*location.Options)) (*location.ListMapsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.ListMapsOutput), args.Error(1)
}

func (m *mockLocationServiceClient) DeleteMap(ctx context.Context,
	params *location.DeleteMapInput,
	_ ...func(*location.Options)) (*location.DeleteMapOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.DeleteMapOutput), args.Error(1)
}

func (m *mockLocationServiceClient) ListRouteCalculators(ctx context.Context,
	params *location.ListRouteCalculatorsInput,
	_ ...func(*location.Options)) (*location.ListRouteCalculatorsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.ListRouteCalculatorsOutput), args.Error(1)
}

func (m *mockLocationServiceClient) DeleteRouteCalculator(ctx context.Context,
	params *location.DeleteRouteCalculatorInput,
	_ ...func(*location.Options)) (*location.DeleteRouteCalculatorOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.DeleteRouteCalculatorOutput), args.Error(1)
}

func (m *mockLocationServiceClient) ListKeys(ctx context.Context,
	params *location.ListKeysInput,
	_ ...func(*location.Options)) (*location.ListKeysOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.ListKeysOutput), args.Error(1)
}

func (m *mockLocationServiceClient) DeleteKey(ctx context.Context,
	params *location.DeleteKeyInput,
	_ ...func(*location.Options)) (*location.DeleteKeyOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.DeleteKeyOutput), args.Error(1)
}

func (m *mockLocationServiceClient) ListPlaceIndexes(
	ctx context.Context, params *location.ListPlaceIndexesInput,
	_ ...func(*location.Options),
) (*location.ListPlaceIndexesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.ListPlaceIndexesOutput), args.Error(1)
}

func (m *mockLocationServiceClient) DeletePlaceIndex(
	ctx context.Context, params *location.DeletePlaceIndexInput,
	_ ...func(*location.Options),
) (*location.DeletePlaceIndexOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*location.DeletePlaceIndexOutput), args.Error(1)
}

var testLocationServiceListerOpts = &nuke.ListerOpts{}
