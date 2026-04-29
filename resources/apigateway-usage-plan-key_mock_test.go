package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	apigatewaytypes "github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

func Test_Mock_APIGatewayUsagePlanKey_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApigatewayClient)

	mockClient.On("GetUsagePlans", mock.Anything, mock.Anything).
		Return(&apigateway.GetUsagePlansOutput{
			Items: []apigatewaytypes.UsagePlan{
				{Id: ptr.String("test-usageplanid")},
			},
		}, nil)

	mockClient.On("GetUsagePlanKeys", mock.Anything, mock.Anything).
		Return(&apigateway.GetUsagePlanKeysOutput{
			Items: []apigatewaytypes.UsagePlanKey{
				{Id: ptr.String("test-keyid"), Name: ptr.String("test-keyname"), Type: ptr.String("test-keytype")},
			},
		}, nil)

	lister := &APIGatewayUsagePlanKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testApigatewayListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*APIGatewayUsagePlanKey)
	a.Equal("test-usageplanid", *r.UsagePlanID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayUsagePlanKey_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApigatewayClient)

	mockClient.On("GetUsagePlans", mock.Anything, mock.Anything).
		Return(&apigateway.GetUsagePlansOutput{
			Items: []apigatewaytypes.UsagePlan{},
		}, nil)

	lister := &APIGatewayUsagePlanKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testApigatewayListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayUsagePlanKey_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApigatewayClient)

	r := &APIGatewayUsagePlanKey{
		svc:         mockClient,
		UsagePlanID: ptr.String("test-usageplanid"),
		KeyID:       ptr.String("test-keyid"),
	}

	mockClient.On("DeleteUsagePlanKey", mock.Anything,
		&apigateway.DeleteUsagePlanKeyInput{
			UsagePlanId: r.UsagePlanID,
			KeyId:       r.KeyID,
		}).Return(&apigateway.DeleteUsagePlanKeyOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_APIGatewayUsagePlanKey_Properties(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayUsagePlanKey{
		UsagePlanID: ptr.String("test-usageplanid"),
		KeyID:       ptr.String("test-keyid"),
		KeyName:     ptr.String("test-keyname"),
		KeyType:     ptr.String("test-keytype"),
	}
	props := r.Properties()
	a.Equal("test-usageplanid", props.Get("UsagePlanId"))
	a.Equal("test-keyid", props.Get("KeyId"))
	a.Equal("test-keyname", props.Get("KeyName"))
	a.Equal("test-keytype", props.Get("KeyType"))
}

func Test_Mock_APIGatewayUsagePlanKey_String(t *testing.T) {
	a := assert.New(t)
	r := &APIGatewayUsagePlanKey{
		KeyID: ptr.String("test-keyid"),
	}
	a.Equal("test-keyid", r.String())
}
