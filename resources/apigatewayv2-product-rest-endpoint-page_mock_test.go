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

func Test_Mock_APIGatewayV2ProductRestEndpointPage_List(t *testing.T) {
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

	mockClient.On("ListProductRestEndpointPages", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListProductRestEndpointPagesOutput{
			Items: []apigatewayv2types.ProductRestEndpointPageSummaryNoBody{
				{
					ProductRestEndpointPageId: ptr.String("rep-789"),
					Endpoint:                  ptr.String("GET /pets"),
				},
			},
		}, nil)

	lister := &APIGatewayV2ProductRestEndpointPageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayV2ProductRestEndpointPage)
	a.Equal("product-123", *r.PortalProductID)
	a.Equal("rep-789", *r.ProductRestEndpointPageID)
	a.Equal("GET /pets", *r.Endpoint)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2ProductRestEndpointPage_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortalProducts", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalProductsOutput{
			Items: []apigatewayv2types.PortalProductSummary{},
		}, nil)

	lister := &APIGatewayV2ProductRestEndpointPageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2ProductRestEndpointPage_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	r := &APIGatewayV2ProductRestEndpointPage{
		svc:                       mockClient,
		PortalProductID:           ptr.String("product-123"),
		ProductRestEndpointPageID: ptr.String("rep-789"),
	}

	mockClient.On("DeleteProductRestEndpointPage", mock.Anything,
		&apigatewayv2.DeleteProductRestEndpointPageInput{
			PortalProductId:           r.PortalProductID,
			ProductRestEndpointPageId: r.ProductRestEndpointPageID,
		}).Return(&apigatewayv2.DeleteProductRestEndpointPageOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2ProductRestEndpointPage_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2ProductRestEndpointPage{
		PortalProductID:           ptr.String("product-123"),
		ProductRestEndpointPageID: ptr.String("rep-789"),
		Endpoint:                  ptr.String("GET /pets"),
	}
	props := r.Properties()
	a.Equal("product-123", props.Get("PortalProductID"))
	a.Equal("rep-789", props.Get("ProductRestEndpointPageID"))
	a.Equal("GET /pets", props.Get("Endpoint"))
}

func Test_Mock_APIGatewayV2ProductRestEndpointPage_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2ProductRestEndpointPage{
		PortalProductID:           ptr.String("product-123"),
		ProductRestEndpointPageID: ptr.String("rep-789"),
		Endpoint:                  ptr.String("GET /pets"),
	}
	a.Equal("product-123 -> rep-789 (GET /pets)", r.String())
}
