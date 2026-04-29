package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connectcases"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockConnectCasesClient struct {
	mock.Mock
}

func (m *mockConnectCasesClient) ListDomains(ctx context.Context,
	params *connectcases.ListDomainsInput,
	_ ...func(*connectcases.Options)) (*connectcases.ListDomainsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connectcases.ListDomainsOutput), args.Error(1)
}

func (m *mockConnectCasesClient) DeleteDomain(ctx context.Context,
	params *connectcases.DeleteDomainInput,
	_ ...func(*connectcases.Options)) (*connectcases.DeleteDomainOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connectcases.DeleteDomainOutput), args.Error(1)
}

func (m *mockConnectCasesClient) ListFields(ctx context.Context,
	params *connectcases.ListFieldsInput,
	_ ...func(*connectcases.Options)) (*connectcases.ListFieldsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connectcases.ListFieldsOutput), args.Error(1)
}

func (m *mockConnectCasesClient) DeleteField(ctx context.Context,
	params *connectcases.DeleteFieldInput,
	_ ...func(*connectcases.Options)) (*connectcases.DeleteFieldOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connectcases.DeleteFieldOutput), args.Error(1)
}

func (m *mockConnectCasesClient) ListTemplates(ctx context.Context,
	params *connectcases.ListTemplatesInput,
	_ ...func(*connectcases.Options)) (*connectcases.ListTemplatesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connectcases.ListTemplatesOutput), args.Error(1)
}

func (m *mockConnectCasesClient) DeleteTemplate(ctx context.Context,
	params *connectcases.DeleteTemplateInput,
	_ ...func(*connectcases.Options)) (*connectcases.DeleteTemplateOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*connectcases.DeleteTemplateOutput), args.Error(1)
}

var testConnectCasesListerOpts = &nuke.ListerOpts{}
