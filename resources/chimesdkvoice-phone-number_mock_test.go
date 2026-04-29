package resources

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkvoice"
	chimesdkvoicetypes "github.com/aws/aws-sdk-go-v2/service/chimesdkvoice/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_ChimeSDKVoicePhoneNumber_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChimeSDKVoiceClient)

	mockClient.On("ListPhoneNumbers", mock.Anything, mock.Anything).
		Return(&chimesdkvoice.ListPhoneNumbersOutput{
			PhoneNumbers: []chimesdkvoicetypes.PhoneNumber{
				{
					PhoneNumberId:   ptr.String("pn-123"),
					E164PhoneNumber: ptr.String("+15551234567"),
					Status:          chimesdkvoicetypes.PhoneNumberStatusAssigned,
					ProductType:     chimesdkvoicetypes.PhoneNumberProductTypeVoiceConnector,
				},
				{
					PhoneNumberId:   ptr.String("pn-456"),
					E164PhoneNumber: ptr.String("+15559876543"),
					Status:          chimesdkvoicetypes.PhoneNumberStatusUnassigned,
				},
			},
		}, nil)

	lister := &ChimeSDKVoicePhoneNumberLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 2)

	r := resources[0].(*ChimeSDKVoicePhoneNumber)
	a.Equal("pn-123", *r.PhoneNumberID)
	a.Equal("+15551234567", *r.E164PhoneNumber)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ChimeSDKVoicePhoneNumber_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChimeSDKVoiceClient)

	mockClient.On("ListPhoneNumbers", mock.Anything, mock.Anything).
		Return(&chimesdkvoice.ListPhoneNumbersOutput{
			PhoneNumbers: []chimesdkvoicetypes.PhoneNumber{},
		}, nil)

	lister := &ChimeSDKVoicePhoneNumberLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ChimeSDKVoicePhoneNumber_Filter_DeleteInProgress(t *testing.T) {
	a := assert.New(t)

	r := &ChimeSDKVoicePhoneNumber{
		PhoneNumberID: ptr.String("pn-123"),
		Status:        chimesdkvoicetypes.PhoneNumberStatusDeleteInProgress,
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already being deleted")
}

func Test_Mock_ChimeSDKVoicePhoneNumber_Filter_Canceled(t *testing.T) {
	a := assert.New(t)

	r := &ChimeSDKVoicePhoneNumber{
		PhoneNumberID: ptr.String("pn-123"),
		Status:        chimesdkvoicetypes.PhoneNumberStatusCancelled,
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "canceled")
}

func Test_Mock_ChimeSDKVoicePhoneNumber_Filter_ReleaseInProgress(t *testing.T) {
	a := assert.New(t)

	r := &ChimeSDKVoicePhoneNumber{
		PhoneNumberID: ptr.String("pn-123"),
		Status:        chimesdkvoicetypes.PhoneNumberStatusReleaseInProgress,
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "release in progress")
}

func Test_Mock_ChimeSDKVoicePhoneNumber_Filter_Active(t *testing.T) {
	a := assert.New(t)

	r := &ChimeSDKVoicePhoneNumber{
		PhoneNumberID: ptr.String("pn-123"),
		Status:        chimesdkvoicetypes.PhoneNumberStatusAssigned,
	}

	err := r.Filter()
	a.Nil(err)
}

func Test_Mock_ChimeSDKVoicePhoneNumber_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChimeSDKVoiceClient)

	r := &ChimeSDKVoicePhoneNumber{
		svc:           mockClient,
		PhoneNumberID: ptr.String("pn-123"),
	}

	mockClient.On("DeletePhoneNumber", mock.Anything, &chimesdkvoice.DeletePhoneNumberInput{
		PhoneNumberId: r.PhoneNumberID,
	}).Return(&chimesdkvoice.DeletePhoneNumberOutput{}, nil)

	err := r.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ChimeSDKVoicePhoneNumber_Properties(t *testing.T) {
	a := assert.New(t)

	r := ChimeSDKVoicePhoneNumber{
		PhoneNumberID:   ptr.String("pn-123"),
		E164PhoneNumber: ptr.String("+15551234567"),
		Status:          chimesdkvoicetypes.PhoneNumberStatusAssigned,
	}

	props := r.Properties()
	a.Equal("pn-123", props.Get("PhoneNumberID"))
	a.Equal("+15551234567", props.Get("E164PhoneNumber"))
}

func Test_Mock_ChimeSDKVoicePhoneNumber_String(t *testing.T) {
	a := assert.New(t)
	r := ChimeSDKVoicePhoneNumber{
		PhoneNumberID:   ptr.String("pn-123"),
		E164PhoneNumber: ptr.String("+15551234567"),
	}
	a.Equal("+15551234567", r.String())
}

func Test_Mock_ChimeSDKVoicePhoneNumber_String_NoE164(t *testing.T) {
	a := assert.New(t)
	r := ChimeSDKVoicePhoneNumber{
		PhoneNumberID: ptr.String("pn-123"),
	}
	a.Equal("pn-123", r.String())
}
