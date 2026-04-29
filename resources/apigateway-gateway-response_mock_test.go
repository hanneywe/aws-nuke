package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func Test_Mock_APIGatewayGatewayResponse_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2Client)

	mockClient.On("GetRestApis", mock.Anything, mock.Anything).
		Return(&apigateway.GetRestApisOutput{
			Items: []apigatewaytypes.RestApi{
				{Id: ptr.String("api-123")},
			},
		}, nil)

	mockClient.On("GetGatewayResponses", mock.Anything, mock.Anything).
		Return(&apigateway.GetGatewayResponsesOutput{
			Items: []apigatewaytypes.GatewayResponse{
				{
					ResponseType:    apigatewaytypes.GatewayResponseTypeDefault4xx,
					StatusCode:      ptr.String("400"),
					DefaultResponse: false,
				},
			},
		}, nil)

	lister := &APIGatewayGatewayResponseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayGatewayResponse)
	a.Equal("api-123", *r.RestAPIID)
	a.Equal(apigatewaytypes.GatewayResponseTypeDefault4xx, r.ResponseType)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayGatewayResponse_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2Client)

	mockClient.On("GetRestApis", mock.Anything, mock.Anything).
		Return(&apigateway.GetRestApisOutput{
			Items: []apigatewaytypes.RestApi{},
		}, nil)

	lister := &APIGatewayGatewayResponseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayGatewayResponse_List_MultipleAPIs(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2Client)

	mockClient.On("GetRestApis", mock.Anything, mock.Anything).
		Return(&apigateway.GetRestApisOutput{
			Items: []apigatewaytypes.RestApi{
				{Id: ptr.String("api-1")},
				{Id: ptr.String("api-2")},
			},
		}, nil)

	mockClient.On("GetGatewayResponses", mock.Anything,
		&apigateway.GetGatewayResponsesInput{
			RestApiId: ptr.String("api-1"),
		}).Return(&apigateway.GetGatewayResponsesOutput{
		Items: []apigatewaytypes.GatewayResponse{
			{
				ResponseType: apigatewaytypes.GatewayResponseTypeDefault4xx,
				StatusCode:   ptr.String("400"),
			},
		},
	}, nil)

	mockClient.On("GetGatewayResponses", mock.Anything,
		&apigateway.GetGatewayResponsesInput{
			RestApiId: ptr.String("api-2"),
		}).Return(&apigateway.GetGatewayResponsesOutput{
		Items: []apigatewaytypes.GatewayResponse{
			{
				ResponseType: apigatewaytypes.GatewayResponseTypeDefault5xx,
				StatusCode:   ptr.String("500"),
			},
		},
	}, nil)

	lister := &APIGatewayGatewayResponseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAPIGatewayV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 2)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayGatewayResponse_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAPIGatewayV2Client)

	r := &APIGatewayGatewayResponse{
		svc:          mockClient,
		RestAPIID:    ptr.String("api-123"),
		ResponseType: apigatewaytypes.GatewayResponseTypeDefault4xx,
	}

	mockClient.On("DeleteGatewayResponse", mock.Anything,
		&apigateway.DeleteGatewayResponseInput{
			RestApiId:    r.RestAPIID,
			ResponseType: r.ResponseType,
		}).Return(&apigateway.DeleteGatewayResponseOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayGatewayResponse_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayGatewayResponse{
		RestAPIID:    ptr.String("api-123"),
		ResponseType: apigatewaytypes.GatewayResponseTypeDefault4xx,
		StatusCode:   ptr.String("400"),
	}
	props := r.Properties()
	a.Equal("api-123", props.Get("RestAPIID"))
	a.Equal("400", props.Get("StatusCode"))
}

func Test_Mock_APIGatewayGatewayResponse_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayGatewayResponse{
		RestAPIID:    ptr.String("api-123"),
		ResponseType: apigatewaytypes.GatewayResponseTypeDefault4xx,
	}
	expected := fmt.Sprintf("api-123 -> %s",
		string(apigatewaytypes.GatewayResponseTypeDefault4xx))
	a.Equal(expected, r.String())
}
