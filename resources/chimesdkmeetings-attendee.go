package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkmeetings"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ChimeSDKMeetingsAttendeeResource = "ChimeSDKMeetingsAttendee"

func init() {
	registry.Register(&registry.Registration{
		Name:     ChimeSDKMeetingsAttendeeResource,
		Scope:    nuke.Account,
		Resource: &ChimeSDKMeetingsAttendee{},
		Lister:   &ChimeSDKMeetingsAttendeeLister{},
	})
}

type ChimeSDKMeetingsAttendeeLister struct {
	svc        ChimeSDKMeetingsClient
	MeetingIDs []string
}

func (l *ChimeSDKMeetingsAttendeeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = chimesdkmeetings.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	for _, meetingID := range l.MeetingIDs {
		params := &chimesdkmeetings.ListAttendeesInput{
			MeetingId: &meetingID,
		}

		for {
			resp, err := svc.ListAttendees(ctx, params)
			if err != nil {
				return nil, err
			}

			for i := range resp.Attendees {
				resources = append(resources, &ChimeSDKMeetingsAttendee{
					svc:            svc,
					MeetingID:      &meetingID,
					AttendeeID:     resp.Attendees[i].AttendeeId,
					ExternalUserID: resp.Attendees[i].ExternalUserId,
				})
			}

			if resp.NextToken == nil {
				break
			}
			params.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

type ChimeSDKMeetingsAttendee struct {
	svc            ChimeSDKMeetingsClient
	MeetingID      *string
	AttendeeID     *string
	ExternalUserID *string
}

func (r *ChimeSDKMeetingsAttendee) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAttendee(ctx, &chimesdkmeetings.DeleteAttendeeInput{
		MeetingId:  r.MeetingID,
		AttendeeId: r.AttendeeID,
	})
	return err
}

func (r *ChimeSDKMeetingsAttendee) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ChimeSDKMeetingsAttendee) String() string {
	return *r.AttendeeID
}
