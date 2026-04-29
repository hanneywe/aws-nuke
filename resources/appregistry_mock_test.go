package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAppRegistryClient struct {
	mock.Mock
}

func (m *mockAppRegistryClient) ListAttributeGroups(
	ctx context.Context,
	params *servicecatalogappregistry.ListAttributeGroupsInput,
	_ ...func(*servicecatalogappregistry.Options),
) (*servicecatalogappregistry.ListAttributeGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*servicecatalogappregistry.ListAttributeGroupsOutput), args.Error(1)
}

func (m *mockAppRegistryClient) DeleteAttributeGroup(
	ctx context.Context,
	params *servicecatalogappregistry.DeleteAttributeGroupInput,
	_ ...func(*servicecatalogappregistry.Options),
) (*servicecatalogappregistry.DeleteAttributeGroupOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*servicecatalogappregistry.DeleteAttributeGroupOutput), args.Error(1)
}

var testAppRegistryListerOpts = &nuke.ListerOpts{}
