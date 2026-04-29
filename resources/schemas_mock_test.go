package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/schemas"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSchemasClient struct {
	mock.Mock
}

func (m *mockSchemasClient) ListRegistries(
	ctx context.Context, params *schemas.ListRegistriesInput,
	_ ...func(*schemas.Options),
) (*schemas.ListRegistriesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*schemas.ListRegistriesOutput), args.Error(1)
}

func (m *mockSchemasClient) DeleteRegistry(
	ctx context.Context, params *schemas.DeleteRegistryInput,
	_ ...func(*schemas.Options),
) (*schemas.DeleteRegistryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*schemas.DeleteRegistryOutput), args.Error(1)
}

func (m *mockSchemasClient) ListDiscoverers(
	ctx context.Context, params *schemas.ListDiscoverersInput,
	_ ...func(*schemas.Options),
) (*schemas.ListDiscoverersOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*schemas.ListDiscoverersOutput), args.Error(1)
}

func (m *mockSchemasClient) DeleteDiscoverer(
	ctx context.Context, params *schemas.DeleteDiscovererInput,
	_ ...func(*schemas.Options),
) (*schemas.DeleteDiscovererOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*schemas.DeleteDiscovererOutput), args.Error(1)
}

var testSchemasListerOpts = &nuke.ListerOpts{}
