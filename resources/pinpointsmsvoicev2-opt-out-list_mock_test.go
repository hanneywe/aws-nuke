package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	pinpointtypes "github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2/types"
)

func Test_Mock_PinpointSMSVoiceV2OptOutList_List_One(t *testing.T) {
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

	lister := &PinpointSMSVoiceV2OptOutListLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ol := resources[0].(*PinpointSMSVoiceV2OptOutList)
	a.Equal("my-opt-out-list", *ol.OptOutListName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2OptOutList_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	mockClient.On("DescribeOptOutLists", mock.Anything, mock.Anything).
		Return(&pinpointsmsvoicev2.DescribeOptOutListsOutput{
			OptOutLists: []pinpointtypes.OptOutListInformation{},
		}, nil)

	lister := &PinpointSMSVoiceV2OptOutListLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testPinpointSMSVoiceV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2OptOutList_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockPinpointSMSVoiceV2Client)

	ol := &PinpointSMSVoiceV2OptOutList{
		svc:            mockClient,
		OptOutListName: ptr.String("my-opt-out-list"),
	}

	mockClient.On("DeleteOptOutList", mock.Anything, &pinpointsmsvoicev2.DeleteOptOutListInput{
		OptOutListName: ol.OptOutListName,
	}).Return(&pinpointsmsvoicev2.DeleteOptOutListOutput{}, nil)

	a.NoError(ol.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_PinpointSMSVoiceV2OptOutList_Filter_Default(t *testing.T) {
	a := assert.New(t)
	ol := PinpointSMSVoiceV2OptOutList{OptOutListName: ptr.String("Default")}
	a.Error(ol.Filter())
}

func Test_Mock_PinpointSMSVoiceV2OptOutList_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	ol := PinpointSMSVoiceV2OptOutList{OptOutListName: ptr.String("my-opt-out-list")}
	a.NoError(ol.Filter())
}

func Test_Mock_PinpointSMSVoiceV2OptOutList_Properties(t *testing.T) {
	a := assert.New(t)

	ol := PinpointSMSVoiceV2OptOutList{
		OptOutListName: ptr.String("my-opt-out-list"),
		OptOutListArn:  ptr.String("arn:aws:sms-voice:us-east-1:123456789012:opt-out-list/my-opt-out-list"),
	}

	props := ol.Properties()
	a.Equal("my-opt-out-list", props.Get("OptOutListName"))
	a.Equal("arn:aws:sms-voice:us-east-1:123456789012:opt-out-list/my-opt-out-list", props.Get("OptOutListArn"))
}

func Test_Mock_PinpointSMSVoiceV2OptOutList_String(t *testing.T) {
	a := assert.New(t)
	ol := PinpointSMSVoiceV2OptOutList{OptOutListName: ptr.String("my-opt-out-list")}
	a.Equal("my-opt-out-list", ol.String())
}
