package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkmeetings"
	chimesdkmeetingstypes "github.com/aws/aws-sdk-go-v2/service/chimesdkmeetings/types"
)

func Test_Mock_ChimeSDKMeetingsAttendee_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChimeSDKMeetingsClient)

	mockClient.On("ListAttendees", mock.Anything, mock.Anything).
		Return(&chimesdkmeetings.ListAttendeesOutput{
			Attendees: []chimesdkmeetingstypes.Attendee{
				{
					AttendeeId:     ptr.String("attendee-123"),
					ExternalUserId: ptr.String("user-abc"),
				},
			},
		}, nil)

	lister := &ChimeSDKMeetingsAttendeeLister{
		svc:        mockClient,
		MeetingIDs: []string{"meeting-456"},
	}
	resources, err := lister.List(context.TODO(), testChimeSDKMeetingsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	attendee := resources[0].(*ChimeSDKMeetingsAttendee)
	a.Equal("attendee-123", *attendee.AttendeeID)
	a.Equal("meeting-456", *attendee.MeetingID)
	a.Equal("user-abc", *attendee.ExternalUserID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ChimeSDKMeetingsAttendee_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChimeSDKMeetingsClient)

	lister := &ChimeSDKMeetingsAttendeeLister{
		svc:        mockClient,
		MeetingIDs: []string{},
	}
	resources, err := lister.List(context.TODO(), testChimeSDKMeetingsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ChimeSDKMeetingsAttendee_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChimeSDKMeetingsClient)

	attendee := &ChimeSDKMeetingsAttendee{
		svc:        mockClient,
		MeetingID:  ptr.String("meeting-456"),
		AttendeeID: ptr.String("attendee-123"),
	}

	mockClient.On("DeleteAttendee", mock.Anything, &chimesdkmeetings.DeleteAttendeeInput{
		MeetingId:  attendee.MeetingID,
		AttendeeId: attendee.AttendeeID,
	}).Return(&chimesdkmeetings.DeleteAttendeeOutput{}, nil)

	a.NoError(attendee.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ChimeSDKMeetingsAttendee_Properties(t *testing.T) {
	a := assert.New(t)

	attendee := ChimeSDKMeetingsAttendee{
		MeetingID:      ptr.String("meeting-456"),
		AttendeeID:     ptr.String("attendee-123"),
		ExternalUserID: ptr.String("user-abc"),
	}

	props := attendee.Properties()
	a.Equal("meeting-456", props.Get("MeetingID"))
	a.Equal("attendee-123", props.Get("AttendeeID"))
	a.Equal("user-abc", props.Get("ExternalUserID"))
}

func Test_Mock_ChimeSDKMeetingsAttendee_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("attendee-123", (&ChimeSDKMeetingsAttendee{AttendeeID: ptr.String("attendee-123")}).String())
}
