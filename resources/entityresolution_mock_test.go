package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/entityresolution"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockEntityresolutionClient struct {
	mock.Mock
}

func (m *mockEntityresolutionClient) ListSchemaMappings(
	ctx context.Context, params *entityresolution.ListSchemaMappingsInput,
	_ ...func(*entityresolution.Options),
) (*entityresolution.ListSchemaMappingsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*entityresolution.ListSchemaMappingsOutput), args.Error(1)
}

func (m *mockEntityresolutionClient) DeleteSchemaMapping(
	ctx context.Context, params *entityresolution.DeleteSchemaMappingInput,
	_ ...func(*entityresolution.Options),
) (*entityresolution.DeleteSchemaMappingOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*entityresolution.DeleteSchemaMappingOutput), args.Error(1)
}

var testEntityresolutionListerOpts = &nuke.ListerOpts{}
