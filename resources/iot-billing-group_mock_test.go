package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func Test_Mock_IoTBillingGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListBillingGroups", mock.Anything, mock.Anything).
		Return(&iot.ListBillingGroupsOutput{
			BillingGroups: []iottypes.GroupNameAndArn{
				{
					GroupName: ptr.String("my-billing-group"),
					GroupArn:  ptr.String("arn:aws:iot:us-east-1:123456789012:billinggroup/my-billing-group"),
				},
			},
		}, nil)

	lister := &IoTBillingGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	bg := resources[0].(*IoTBillingGroup)
	a.Equal("my-billing-group", *bg.BillingGroupName)
	a.Equal("arn:aws:iot:us-east-1:123456789012:billinggroup/my-billing-group", *bg.BillingGroupArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTBillingGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListBillingGroups", mock.Anything, mock.Anything).
		Return(&iot.ListBillingGroupsOutput{
			BillingGroups: []iottypes.GroupNameAndArn{},
		}, nil)

	lister := &IoTBillingGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTBillingGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	bg := &IoTBillingGroup{
		svc:              mockClient,
		BillingGroupName: ptr.String("my-billing-group"),
	}

	mockClient.On("DeleteBillingGroup", mock.Anything, &iot.DeleteBillingGroupInput{
		BillingGroupName: bg.BillingGroupName,
	}).Return(&iot.DeleteBillingGroupOutput{}, nil)

	a.NoError(bg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTBillingGroup_Properties(t *testing.T) {
	a := assert.New(t)

	bg := IoTBillingGroup{
		BillingGroupName: ptr.String("my-billing-group"),
		BillingGroupArn:  ptr.String("arn:aws:iot:us-east-1:123456789012:billinggroup/my-billing-group"),
	}

	props := bg.Properties()
	a.Equal("my-billing-group", props.Get("BillingGroupName"))
	a.Equal("arn:aws:iot:us-east-1:123456789012:billinggroup/my-billing-group", props.Get("BillingGroupArn"))
}

func Test_Mock_IoTBillingGroup_String(t *testing.T) {
	a := assert.New(t)
	bg := IoTBillingGroup{BillingGroupName: ptr.String("my-billing-group")}
	a.Equal("my-billing-group", bg.String())
}
