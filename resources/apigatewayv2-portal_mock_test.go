package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAPIGatewayV2PortalClient struct {
	mock.Mock
}

func (m *mockAPIGatewayV2PortalClient) ListPortals(
	ctx context.Context, params *apigatewayv2.ListPortalsInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.ListPortalsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.ListPortalsOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) DeletePortal(
	ctx context.Context, params *apigatewayv2.DeletePortalInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.DeletePortalOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.DeletePortalOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) DisablePortal(
	ctx context.Context, params *apigatewayv2.DisablePortalInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.DisablePortalOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.DisablePortalOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) ListPortalProducts(
	ctx context.Context, params *apigatewayv2.ListPortalProductsInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.ListPortalProductsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.ListPortalProductsOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) DeletePortalProduct(
	ctx context.Context, params *apigatewayv2.DeletePortalProductInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.DeletePortalProductOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.DeletePortalProductOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) ListProductPages(
	ctx context.Context, params *apigatewayv2.ListProductPagesInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.ListProductPagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.ListProductPagesOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) DeleteProductPage(
	ctx context.Context, params *apigatewayv2.DeleteProductPageInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.DeleteProductPageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.DeleteProductPageOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) ListProductRestEndpointPages(
	ctx context.Context, params *apigatewayv2.ListProductRestEndpointPagesInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.ListProductRestEndpointPagesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.ListProductRestEndpointPagesOutput), args.Error(1)
}

func (m *mockAPIGatewayV2PortalClient) DeleteProductRestEndpointPage(
	ctx context.Context, params *apigatewayv2.DeleteProductRestEndpointPageInput,
	_ ...func(*apigatewayv2.Options),
) (*apigatewayv2.DeleteProductRestEndpointPageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*apigatewayv2.DeleteProductRestEndpointPageOutput), args.Error(1)
}

var testAPIGatewayV2PortalListerOpts = &nuke.ListerOpts{}
