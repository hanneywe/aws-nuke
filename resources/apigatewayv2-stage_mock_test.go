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

func Test_Mock_APIGatewayV2Stage_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApigatewayv2Client)

	mockClient.On("GetApis", mock.Anything, mock.Anything).
		Return(&apigatewayv2.GetApisOutput{
			Items: []apigatewayv2types.Api{
				{ApiId: ptr.String("test-apiid")},
			},
		}, nil)

	mockClient.On("GetStages", mock.Anything, mock.Anything).
		Return(&apigatewayv2.GetStagesOutput{
			Items: []apigatewayv2types.Stage{
				{StageName: ptr.String("test-stagename"), Tags: map[string]string{"env": "test"}},
			},
		}, nil)

	lister := &APIGatewayV2StageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testApigatewayv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayV2Stage)
	a.Equal("test-apiid", *r.APIID)
	a.Equal("test-stagename", *r.StageName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Stage_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApigatewayv2Client)

	mockClient.On("GetApis", mock.Anything, mock.Anything).
		Return(&apigatewayv2.GetApisOutput{
			Items: []apigatewayv2types.Api{},
		}, nil)

	lister := &APIGatewayV2StageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testApigatewayv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Stage_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApigatewayv2Client)

	r := &APIGatewayV2Stage{
		svc:       mockClient,
		APIID:     ptr.String("test-apiid"),
		StageName: ptr.String("test-stagename"),
	}

	mockClient.On("DeleteStage", mock.Anything,
		&apigatewayv2.DeleteStageInput{
			ApiId:     r.APIID,
			StageName: r.StageName,
		}).Return(&apigatewayv2.DeleteStageOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayV2Stage_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2Stage{
		APIID:     ptr.String("test-apiid"),
		StageName: ptr.String("test-stagename"),
	}
	props := r.Properties()
	a.Equal("test-apiid", props.Get("ApiId"))
	a.Equal("test-stagename", props.Get("StageName"))
}

func Test_Mock_APIGatewayV2Stage_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayV2Stage{
		StageName: ptr.String("test-stagename"),
	}
	a.Equal("test-stagename", r.String())
}
