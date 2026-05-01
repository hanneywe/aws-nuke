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

func Test_Mock_CleanRoomsCollaboration_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	mockClient.
		On("ListCollaborations", mock.Anything, mock.Anything).
		Return(
			&cleanrooms.ListCollaborationsOutput{
				CollaborationList: []cleanroomstypes.CollaborationSummary{
					{
						Id:   ptr.String("collab-12345"),
						Name: ptr.String("test-collaboration"),
					},
				},
			}, nil,
		)

	lister := &CleanRoomsCollaborationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testCleanRoomsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	collaboration := resources[0].(*CleanRoomsCollaboration)
	assertions.Equal("collab-12345", *collaboration.CollaborationIdentifier)
	assertions.Equal("test-collaboration", *collaboration.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsCollaboration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	mockClient.
		On("ListCollaborations", mock.Anything, mock.Anything).
		Return(
			&cleanrooms.ListCollaborationsOutput{
				CollaborationList: []cleanroomstypes.CollaborationSummary{},
			}, nil,
		)

	lister := &CleanRoomsCollaborationLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testCleanRoomsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsCollaboration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockCleanRoomsClient)

	collaboration := &CleanRoomsCollaboration{
		svc:                     mockClient,
		CollaborationIdentifier: ptr.String("collab-12345"),
		Name:                    ptr.String("test-collaboration"),
	}

	mockClient.
		On(
			"DeleteCollaboration",
			mock.Anything,
			&cleanrooms.DeleteCollaborationInput{
				CollaborationIdentifier: collaboration.CollaborationIdentifier,
			},
		).
		Return(&cleanrooms.DeleteCollaborationOutput{}, nil)

	err := collaboration.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CleanRoomsCollaboration_Properties(t *testing.T) {
	assertions := assert.New(t)

	collaboration := CleanRoomsCollaboration{
		CollaborationIdentifier: ptr.String("collab-12345"),
		Name:                    ptr.String("test-collaboration"),
	}

	properties := collaboration.Properties()

	assertions.Equal("collab-12345", properties.Get("CollaborationIdentifier"))
	assertions.Equal("test-collaboration", properties.Get("Name"))
}

func Test_Mock_CleanRoomsCollaboration_String(t *testing.T) {
	assertions := assert.New(t)

	collaboration := CleanRoomsCollaboration{
		CollaborationIdentifier: ptr.String("collab-12345"),
		Name:                    ptr.String("test-collaboration"),
	}

	assertions.Equal("collab-12345", collaboration.String())
}
