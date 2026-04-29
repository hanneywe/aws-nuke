package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testEventBridgeListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_EventBridgeArchive_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	mockClient.
		On("ListArchives", mock.Anything, mock.Anything).
		Return(
			&eventbridge.ListArchivesOutput{
				Archives: []eventbridgetypes.Archive{
					{
						ArchiveName: ptr.String("test-archive"),
					},
				},
			}, nil,
		)

	lister := &EventBridgeArchiveLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEventBridgeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	archive := resources[0].(*EventBridgeArchive)
	assertions.Equal("test-archive", *archive.ArchiveName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeArchive_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	mockClient.
		On("ListArchives", mock.Anything, mock.Anything).
		Return(
			&eventbridge.ListArchivesOutput{
				Archives: []eventbridgetypes.Archive{},
			}, nil,
		)

	lister := &EventBridgeArchiveLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testEventBridgeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeArchive_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockEventBridgeClient)

	archive := &EventBridgeArchive{
		svc:         mockClient,
		ArchiveName: ptr.String("test-archive"),
	}

	mockClient.
		On(
			"DeleteArchive",
			mock.Anything,
			&eventbridge.DeleteArchiveInput{
				ArchiveName: archive.ArchiveName,
			},
		).
		Return(&eventbridge.DeleteArchiveOutput{}, nil)

	err := archive.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EventBridgeArchive_Properties(t *testing.T) {
	assertions := assert.New(t)

	archive := EventBridgeArchive{
		ArchiveName: ptr.String("test-archive"),
	}

	properties := archive.Properties()

	assertions.Equal("test-archive", properties.Get("ArchiveName"))
}

func Test_Mock_EventBridgeArchive_String(t *testing.T) {
	assertions := assert.New(t)

	archive := EventBridgeArchive{
		ArchiveName: ptr.String("test-archive"),
	}

	assertions.Equal("test-archive", archive.String())
}
