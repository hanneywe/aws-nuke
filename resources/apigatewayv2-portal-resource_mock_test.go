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

func Test_Mock_APIGatewayV2Portal_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortals", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalsOutput{
			Items: []apigatewayv2types.PortalSummary{
				{
					PortalId: ptr.String("portal-123"),
					PortalContent: &apigatewayv2types.PortalContent{
						DisplayName: ptr.String("My Portal"),
					},
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &APIGatewayV2PortalLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayV2Portal)
	a.Equal("portal-123", *r.PortalID)
	a.Equal("My Portal", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Portal_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	mockClient.On("ListPortals", mock.Anything, mock.Anything).
		Return(&apigatewayv2.ListPortalsOutput{
			Items: []apigatewayv2types.PortalSummary{},
		}, nil)

	lister := &APIGatewayV2PortalLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2PortalListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Portal_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	r := &APIGatewayV2Portal{
		svc:           mockClient,
		PortalID:      ptr.String("portal-123"),
		PublishStatus: apigatewayv2types.PublishStatusPublished,
	}

	mockClient.On("DisablePortal", mock.Anything,
		&apigatewayv2.DisablePortalInput{
			PortalId: r.PortalID,
		}).Return(&apigatewayv2.DisablePortalOutput{}, nil)

	err := r.Remove(context.TODO())
	a.Error(err)
	a.Contains(err.Error(), "waiting for portal to be disabled")
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Portal_Remove_AlreadyDisabled(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2PortalClient)

	r := &APIGatewayV2Portal{
		svc:           mockClient,
		PortalID:      ptr.String("portal-123"),
		PublishStatus: apigatewayv2types.PublishStatusDisabled,
	}

	mockClient.On("DeletePortal", mock.Anything,
		&apigatewayv2.DeletePortalInput{
			PortalId: r.PortalID,
		}).Return(&apigatewayv2.DeletePortalOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Portal_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2Portal{
		PortalID: ptr.String("portal-123"),
		Name:     ptr.String("My Portal"),
		Tags:     map[string]string{"env": "test"},
	}
	props := r.Properties()
	a.Equal("portal-123", props.Get("PortalID"))
	a.Equal("My Portal", props.Get("Name"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_APIGatewayV2Portal_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2Portal{
		PortalID: ptr.String("portal-123"),
		Name:     ptr.String("My Portal"),
	}
	a.Equal("portal-123 (My Portal)", r.String())
}

func Test_Mock_APIGatewayV2Portal_String_NoName(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2Portal{
		PortalID: ptr.String("portal-123"),
	}
	a.Equal("portal-123", r.String())
}
