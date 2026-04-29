package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"
)

func Test_Mock_PinpointSMSVoiceV2OptedOutNumber_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeOptOutLists", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeOptOutListsOutput{
			OptOutLists: []pinpointtypes.OptOutListInformation{
				{
					OptOutListName: ptr.String("my-opt-out-list"),
					OptOutListArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:opt-out-list/my-opt-out-list"),
				},
			},
		}, nil)

	now := time.Now()
	mockClient.On("DescribeOptedOutNumbers", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeOptedOutNumbersOutput{
			OptedOutNumbers: []pinpointtypes.OptedOutNumberInformation{
				{
					OptedOutNumber:    ptr.String("+12065551234"),
					OptedOutTimestamp: &now,
					EndUserOptedOut:   true,
				},
			},
		}, nil)

	lister := &PinpointSMSVoiceV2OptedOutNumberLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	num := resources[0].(*PinpointSMSVoiceV2OptedOutNumber)
	a.Equal("+12065551234", *num.OptedOutNumber)
	a.Equal("my-opt-out-list", *num.OptOutListName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2OptedOutNumber_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeOptOutLists", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeOptOutListsOutput{
			OptOutLists: []pinpointtypes.OptOutListInformation{
				{
					OptOutListName: ptr.String("my-opt-out-list"),
					OptOutListArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:opt-out-list/my-opt-out-list"),
				},
			},
		}, nil)

	mockClient.On("DescribeOptedOutNumbers", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeOptedOutNumbersOutput{
			OptedOutNumbers: []pinpointtypes.OptedOutNumberInformation{},
		}, nil)

	lister := &PinpointSMSVoiceV2OptedOutNumberLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2OptedOutNumber_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	num := &PinpointSMSVoiceV2OptedOutNumber{
		svc:            mockClient,
		OptedOutNumber: ptr.String("+12065551234"),
		OptOutListName: ptr.String("my-opt-out-list"),
	}

	mockClient.On("DeleteOptedOutNumber", mock.Anything, &pinpointsmsvoicev2.DeleteOptedOutNumberInput{
		OptOutListName: num.OptOutListName,
		OptedOutNumber: num.OptedOutNumber,
	}).Return(&pinpointsmsvoicev2.DeleteOptedOutNumberOutput{}, nil)

	a.NoError(num.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2OptedOutNumber_Properties(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	num := PinpointSMSVoiceV2OptedOutNumber{
		OptedOutNumber:    ptr.String("+12065551234"),
		OptOutListName:    ptr.String("my-opt-out-list"),
		OptedOutTimestamp: &now,
		EndUserOptedOut:   ptr.Bool(true),
	}

	props := num.Properties()
	a.Equal("+12065551234", props.Get("OptedOutNumber"))
	a.Equal("my-opt-out-list", props.Get("OptOutListName"))
}

func Test_Mock_PinpointSMSVoiceV2OptedOutNumber_String(t *testing.T) {
	a := assert.New(t)
	num := PinpointSMSVoiceV2OptedOutNumber{
		OptedOutNumber: ptr.String("+12065551234"),
		OptOutListName: ptr.String("my-opt-out-list"),
	}
	a.Equal("my-opt-out-list -> +12065551234", num.String())
}
