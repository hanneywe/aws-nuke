package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSESv2Client struct {
	mock.Mock
}

func (m *mockSESv2Client) ListConfigurationSets(
	ctx context.Context, params *sesv2.ListConfigurationSetsInput,
	_ ...func(*sesv2.Options),
) (*sesv2.ListConfigurationSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.ListConfigurationSetsOutput), args.Error(1)
}

func (m *mockSESv2Client) DeleteConfigurationSet(
	ctx context.Context, params *sesv2.DeleteConfigurationSetInput,
	_ ...func(*sesv2.Options),
) (*sesv2.DeleteConfigurationSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.DeleteConfigurationSetOutput), args.Error(1)
}

func (m *mockSESv2Client) ListDedicatedIpPools(
	ctx context.Context, params *sesv2.ListDedicatedIpPoolsInput,
	_ ...func(*sesv2.Options),
) (*sesv2.ListDedicatedIpPoolsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.ListDedicatedIpPoolsOutput), args.Error(1)
}

func (m *mockSESv2Client) DeleteDedicatedIpPool(
	ctx context.Context, params *sesv2.DeleteDedicatedIpPoolInput,
	_ ...func(*sesv2.Options),
) (*sesv2.DeleteDedicatedIpPoolOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.DeleteDedicatedIpPoolOutput), args.Error(1)
}

func (m *mockSESv2Client) ListEmailIdentities(
	ctx context.Context, params *sesv2.ListEmailIdentitiesInput,
	_ ...func(*sesv2.Options),
) (*sesv2.ListEmailIdentitiesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.ListEmailIdentitiesOutput), args.Error(1)
}

func (m *mockSESv2Client) DeleteEmailIdentity(
	ctx context.Context, params *sesv2.DeleteEmailIdentityInput,
	_ ...func(*sesv2.Options),
) (*sesv2.DeleteEmailIdentityOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.DeleteEmailIdentityOutput), args.Error(1)
}

func (m *mockSESv2Client) ListEmailTemplates(
	ctx context.Context, params *sesv2.ListEmailTemplatesInput,
	_ ...func(*sesv2.Options),
) (*sesv2.ListEmailTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.ListEmailTemplatesOutput), args.Error(1)
}

func (m *mockSESv2Client) DeleteEmailTemplate(
	ctx context.Context, params *sesv2.DeleteEmailTemplateInput,
	_ ...func(*sesv2.Options),
) (*sesv2.DeleteEmailTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sesv2.DeleteEmailTemplateOutput), args.Error(1)
}

var testSESv2ListerOpts = &nuke.ListerOpts{}
