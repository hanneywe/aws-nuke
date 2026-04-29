package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	apigatewayv2types "github.com/aws/aws-sdk-go-v2/service/apigatewayv2/types"
)

func Test_Mock_APIGatewayV2ProductPage_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortalProducts", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalProductsOutput{
			Items: []apigatewayv2types.PortalProductSummary{
				{
					PortalProductId: ptr.String("product-123"),
					DisplayName:     ptr.String("Product"),
					Description:     ptr.String("desc"),
				},
			},
		}, nil)

	mockClient.On("ListProductPages", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListProductPagesOutput{
			Items: []apigatewayv2types.ProductPageSummaryNoBody{
				{
					ProductPageId: ptr.String("page-456"),
					PageTitle:     ptr.String("Getting Started"),
				},
			},
		}, nil)

	lister := &APIGatewayV2ProductPageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayV2ProductPage)
	a.Equal("product-123", *r.PortalProductID)
	a.Equal("page-456", *r.ProductPageID)
	a.Equal("Getting Started", *r.PageTitle)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2ProductPage_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortalProducts", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalProductsOutput{
			Items: []apigatewayv2types.PortalProductSummary{},
		}, nil)

	lister := &APIGatewayV2ProductPageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2ProductPage_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	r := &APIGatewayV2ProductPage{
		svc:             mockClient,
		PortalProductID: ptr.String("product-123"),
		ProductPageID:   ptr.String("page-456"),
	}

	mockClient.On("DeleteProductPage", mock.Anything,
		&apigatewayv2.DeleteProductPageInput{
			PortalProductId: r.PortalProductID,
			ProductPageId:   r.ProductPageID,
		}).Return(&apigatewayv2.DeleteProductPageOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2ProductPage_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2ProductPage{
		PortalProductID: ptr.String("product-123"),
		ProductPageID:   ptr.String("page-456"),
		PageTitle:       ptr.String("Getting Started"),
	}
	props := r.Properties()
	a.Equal("product-123", props.Get("PortalProductID"))
	a.Equal("page-456", props.Get("ProductPageID"))
	a.Equal("Getting Started", props.Get("PageTitle"))
}

func Test_Mock_APIGatewayV2ProductPage_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2ProductPage{
		PortalProductID: ptr.String("product-123"),
		ProductPageID:   ptr.String("page-456"),
		PageTitle:       ptr.String("Getting Started"),
	}
	a.Equal("product-123 -> page-456 (Getting Started)", r.String())
}
