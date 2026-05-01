package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"
	cleanroomstypes "github.com/aws/aws-sdk-go-v2/service/cleanrooms/types"
)

func Test_Mock_CleanRoomsMembership_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	mockClient.On("ListMemberships", mock.Anything, mock.Anything).
		Return(&cleanrooms.ListMembershipsOutput{
			MembershipSummaries: []cleanroomstypes.MembershipSummary{
				{
					Id:                ptr.String("membership-12345"),
					CollaborationName: ptr.String("test-collaboration"),
					CollaborationId:   ptr.String("collab-12345"),
				},
			},
		}, nil)

	lister := &CleanRoomsMembershipLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCleanRoomsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CleanRoomsMembership)
	a.Equal("membership-12345", *r.MembershipIdentifier)
	a.Equal("test-collaboration", *r.CollaborationName)
	a.Equal("collab-12345", *r.CollaborationID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsMembership_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	mockClient.On("ListMemberships", mock.Anything, mock.Anything).
		Return(&cleanrooms.ListMembershipsOutput{
			MembershipSummaries: []cleanroomstypes.MembershipSummary{},
		}, nil)

	lister := &CleanRoomsMembershipLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCleanRoomsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsMembership_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	r := &CleanRoomsMembership{
		svc:                  mockClient,
		MembershipIdentifier: ptr.String("membership-12345"),
	}

	mockClient.On("DeleteMembership", mock.Anything,
		&cleanrooms.DeleteMembershipInput{
			MembershipIdentifier: r.MembershipIdentifier,
		}).Return(&cleanrooms.DeleteMembershipOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsMembership_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CleanRoomsMembership{
		MembershipIdentifier: ptr.String("membership-12345"),
		CollaborationName:    ptr.String("test-collaboration"),
		CollaborationID:      ptr.String("collab-12345"),
	}
	props := r.Properties()
	a.Equal("membership-12345", props.Get("MembershipIdentifier"))
	a.Equal("test-collaboration", props.Get("CollaborationName"))
	a.Equal("collab-12345", props.Get("CollaborationID"))
}

func Test_Mock_CleanRoomsMembership_String(t *testing.T) {
	a := assert.New(t)
	r := &CleanRoomsMembership{
		MembershipIdentifier: ptr.String("membership-12345"),
	}
	a.Equal("membership-12345", r.String())
}
