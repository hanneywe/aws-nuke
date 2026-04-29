package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/redshift"
	redshifttypes "github.com/aws/aws-sdk-go-v2/service/redshift/types"
)

func Test_Mock_RedshiftEventSubscription_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeEventSubscriptions", mock.Anything, mock.Anything).
		Return(&redshift.DescribeEventSubscriptionsOutput{
			EventSubscriptionsList: []redshifttypes.EventSubscription{
				{
					CustSubscriptionId: ptr.String("my-subscription"),
				},
			},
		}, nil)

	lister := &RedshiftEventSubscriptionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	subscription := resources[0].(*RedshiftEventSubscription)
	assertions.Equal("my-subscription", *subscription.CustSubscriptionID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftEventSubscription_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRedshiftClient)

	mockClient.On("DescribeEventSubscriptions", mock.Anything, mock.Anything).
		Return(&redshift.DescribeEventSubscriptionsOutput{
			EventSubscriptionsList: []redshifttypes.EventSubscription{},
		}, nil)

	lister := &RedshiftEventSubscriptionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRedshiftListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftEventSubscription_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRedshiftClient)

	subscription := &RedshiftEventSubscription{
		svc:                mockClient,
		CustSubscriptionID: ptr.String("my-subscription"),
	}

	mockClient.On("DeleteEventSubscription", mock.Anything, &redshift.DeleteEventSubscriptionInput{
		SubscriptionName: subscription.CustSubscriptionID,
	}).Return(&redshift.DeleteEventSubscriptionOutput{}, nil)

	err := subscription.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RedshiftEventSubscription_Properties(t *testing.T) {
	assertions := assert.New(t)

	subscription := RedshiftEventSubscription{
		CustSubscriptionID: ptr.String("my-subscription"),
	}

	properties := subscription.Properties()
	assertions.Equal("my-subscription", properties.Get("CustSubscriptionId"))
}

func Test_Mock_RedshiftEventSubscription_String(t *testing.T) {
	assertions := assert.New(t)
	subscription := RedshiftEventSubscription{CustSubscriptionID: ptr.String("my-subscription")}
	assertions.Equal("my-subscription", subscription.String())
}
