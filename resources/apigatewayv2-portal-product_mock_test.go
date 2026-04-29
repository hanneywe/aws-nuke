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

func Test_Mock_APIGatewayV2PortalProduct_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortalProducts", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalProductsOutput{
			Items: []apigatewayv2types.PortalProductSummary{
				{
					PortalProductId: ptr.String("product-123"),
					DisplayName:     ptr.String("My Product"),
					Description:     ptr.String("desc"),
					Tags:            map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &APIGatewayV2PortalProductLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayV2PortalProduct)
	a.Equal("product-123", *r.PortalProductID)
	a.Equal("My Product", *r.DisplayName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2PortalProduct_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortalProducts", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalProductsOutput{
			Items: []apigatewayv2types.PortalProductSummary{},
		}, nil)

	lister := &APIGatewayV2PortalProductLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2PortalProduct_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	r := &APIGatewayV2PortalProduct{
		svc:             mockClient,
		PortalProductID: ptr.String("product-123"),
	}

	mockClient.On("DeletePortalProduct", mock.Anything,
		&apigatewayv2.DeletePortalProductInput{
			PortalProductId: r.PortalProductID,
		}).Return(&apigatewayv2.DeletePortalProductOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2PortalProduct_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2PortalProduct{
		PortalProductID: ptr.String("product-123"),
		DisplayName:     ptr.String("My Product"),
		Tags:            map[string]string{"env": "test"},
	}
	props := r.Properties()
	a.Equal("product-123", props.Get("PortalProductID"))
	a.Equal("My Product", props.Get("DisplayName"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_APIGatewayV2PortalProduct_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2PortalProduct{
		PortalProductID: ptr.String("product-123"),
		DisplayName:     ptr.String("My Product"),
	}
	a.Equal("product-123 (My Product)", r.String())
}
