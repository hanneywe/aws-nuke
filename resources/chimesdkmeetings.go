package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkmeetings"
)

type ChimeSDKMeetingsClient interface {
	ListAttendees(ctx context.Context, params *chimesdkmeetings.ListAttendeesInput,
		optFns ...func(*chimesdkmeetings.Options)) (*chimesdkmeetings.ListAttendeesOutput, error)
	DeleteAttendee(ctx context.Context, params *chimesdkmeetings.DeleteAttendeeInput,
		optFns ...func(*chimesdkmeetings.Options)) (*chimesdkmeetings.DeleteAttendeeOutput, error)
}
