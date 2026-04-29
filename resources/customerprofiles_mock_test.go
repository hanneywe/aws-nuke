package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockCustomerProfilesClient struct {
	mock.Mock
}

func (m *mockCustomerProfilesClient) ListDomains(ctx context.Context,
	params *customerprofiles.ListDomainsInput,
	_ ...func(*customerprofiles.Options)) (*customerprofiles.ListDomainsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*customerprofiles.ListDomainsOutput), args.Error(1)
}

func (m *mockCustomerProfilesClient) DeleteDomain(ctx context.Context,
	params *customerprofiles.DeleteDomainInput,
	_ ...func(*customerprofiles.Options)) (*customerprofiles.DeleteDomainOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*customerprofiles.DeleteDomainOutput), args.Error(1)
}

func (m *mockCustomerProfilesClient) ListSegmentDefinitions(ctx context.Context,
	params *customerprofiles.ListSegmentDefinitionsInput,
	_ ...func(*customerprofiles.Options)) (*customerprofiles.ListSegmentDefinitionsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*customerprofiles.ListSegmentDefinitionsOutput), args.Error(1)
}

func (m *mockCustomerProfilesClient) DeleteSegmentDefinition(ctx context.Context,
	params *customerprofiles.DeleteSegmentDefinitionInput,
	_ ...func(*customerprofiles.Options)) (*customerprofiles.DeleteSegmentDefinitionOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*customerprofiles.DeleteSegmentDefinitionOutput), args.Error(1)
}

func (m *mockCustomerProfilesClient) SearchProfiles(ctx context.Context,
	params *customerprofiles.SearchProfilesInput,
	_ ...func(*customerprofiles.Options)) (*customerprofiles.SearchProfilesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*customerprofiles.SearchProfilesOutput), args.Error(1)
}

func (m *mockCustomerProfilesClient) DeleteProfile(ctx context.Context,
	params *customerprofiles.DeleteProfileInput,
	_ ...func(*customerprofiles.Options)) (*customerprofiles.DeleteProfileOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*customerprofiles.DeleteProfileOutput), args.Error(1)
}

var testCustomerProfilesListerOpts = &nuke.ListerOpts{}
