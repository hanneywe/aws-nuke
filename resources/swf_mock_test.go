package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/swf"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSWFClient struct {
	mock.Mock
}

func (m *mockSWFClient) ListDomains(
	ctx context.Context, params *swf.ListDomainsInput,
	_ ...func(*swf.Options),
) (*swf.ListDomainsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*swf.ListDomainsOutput), args.Error(1)
}

func (m *mockSWFClient) DeprecateDomain(
	ctx context.Context, params *swf.DeprecateDomainInput,
	_ ...func(*swf.Options),
) (*swf.DeprecateDomainOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*swf.DeprecateDomainOutput), args.Error(1)
}

func (m *mockSWFClient) ListActivityTypes(
	ctx context.Context, params *swf.ListActivityTypesInput,
	_ ...func(*swf.Options),
) (*swf.ListActivityTypesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*swf.ListActivityTypesOutput), args.Error(1)
}

func (m *mockSWFClient) DeprecateActivityType(
	ctx context.Context, params *swf.DeprecateActivityTypeInput,
	_ ...func(*swf.Options),
) (*swf.DeprecateActivityTypeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*swf.DeprecateActivityTypeOutput), args.Error(1)
}

var testSWFListerOpts = &nuke.ListerOpts{}
