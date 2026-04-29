package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkmeetings"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockChimeSDKMeetingsClient struct {
	mock.Mock
}

func (m *mockChimeSDKMeetingsClient) ListAttendees(
	ctx context.Context, params *chimesdkmeetings.ListAttendeesInput, _ ...func(*chimesdkmeetings.Options),
) (*chimesdkmeetings.ListAttendeesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chimesdkmeetings.ListAttendeesOutput), args.Error(1)
}

func (m *mockChimeSDKMeetingsClient) DeleteAttendee(
	ctx context.Context, params *chimesdkmeetings.DeleteAttendeeInput, _ ...func(*chimesdkmeetings.Options),
) (*chimesdkmeetings.DeleteAttendeeOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*chimesdkmeetings.DeleteAttendeeOutput), args.Error(1)
}

var testChimeSDKMeetingsListerOpts = &nuke.ListerOpts{}
