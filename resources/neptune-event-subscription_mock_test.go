package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/neptune"
	neptunetypes "github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

func Test_Mock_NeptuneEventSubscription_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeEventSubscriptions", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeEventSubscriptionsOutput{
				EventSubscriptionsList: []neptunetypes.EventSubscription{
					{
						CustSubscriptionId:   ptr.String("my-neptune-subscription"),
						EventSubscriptionArn: ptr.String("arn:aws:rds:us-east-1:123456789012:es:my-neptune-subscription"),
					},
				},
			}, nil,
		)

	lister := &NeptuneEventSubscriptionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	eventSubscription := resources[0].(*NeptuneEventSubscription)
	assertions.Equal("my-neptune-subscription", *eventSubscription.CustSubscriptionID)
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:es:my-neptune-subscription", *eventSubscription.EventSubscriptionArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneEventSubscription_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	mockClient.
		On("DescribeEventSubscriptions", mock.Anything, mock.Anything).
		Return(
			&neptune.DescribeEventSubscriptionsOutput{
				EventSubscriptionsList: []neptunetypes.EventSubscription{},
			}, nil,
		)

	lister := &NeptuneEventSubscriptionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testNeptuneV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneEventSubscription_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockNeptuneV2Client)

	eventSubscription := &NeptuneEventSubscription{
		svc:                  mockClient,
		CustSubscriptionID:   ptr.String("my-neptune-subscription"),
		EventSubscriptionArn: ptr.String("arn:aws:rds:us-east-1:123456789012:es:my-neptune-subscription"),
	}

	mockClient.
		On("DeleteEventSubscription", mock.Anything,
			&neptune.DeleteEventSubscriptionInput{
				SubscriptionName: eventSubscription.CustSubscriptionID,
			},
		).
		Return(&neptune.DeleteEventSubscriptionOutput{}, nil)

	err := eventSubscription.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneEventSubscription_Properties(t *testing.T) {
	assertions := assert.New(t)

	eventSubscription := NeptuneEventSubscription{
		CustSubscriptionID:   ptr.String("my-neptune-subscription"),
		EventSubscriptionArn: ptr.String("arn:aws:rds:us-east-1:123456789012:es:my-neptune-subscription"),
	}

	properties := eventSubscription.Properties()

	assertions.Equal("my-neptune-subscription", properties.Get("CustSubscriptionId"))
	assertions.Equal("arn:aws:rds:us-east-1:123456789012:es:my-neptune-subscription", properties.Get("EventSubscriptionArn"))
}

func Test_Mock_NeptuneEventSubscription_String(t *testing.T) {
	assertions := assert.New(t)

	eventSubscription := NeptuneEventSubscription{
		CustSubscriptionID: ptr.String("my-neptune-subscription"),
	}

	assertions.Equal("my-neptune-subscription", eventSubscription.String())
}
