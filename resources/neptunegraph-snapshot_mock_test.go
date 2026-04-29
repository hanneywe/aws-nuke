package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
	neptunegraphtypes "github.com/aws/aws-sdk-go-v2/service/neptunegraph/types"
)

func Test_Mock_NeptuneGraphSnapshot_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNeptuneGraphClient)

	mockClient.On("ListGraphSnapshots", mock.Anything, mock.Anything).
		Return(&neptunegraph.ListGraphSnapshotsOutput{
			GraphSnapshots: []neptunegraphtypes.GraphSnapshotSummary{
				{
					Id:     ptr.String("gs-abc123"),
					Name:   ptr.String("my-snapshot"),
					Status: neptunegraphtypes.SnapshotStatusAvailable,
				},
			},
		}, nil)

	lister := &NeptuneGraphSnapshotLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNeptuneGraphListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	snap := resources[0].(*NeptuneGraphSnapshot)
	a.Equal("gs-abc123", *snap.SnapshotID)
	a.Equal("my-snapshot", *snap.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGraphSnapshot_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNeptuneGraphClient)

	mockClient.On("ListGraphSnapshots", mock.Anything, mock.Anything).
		Return(&neptunegraph.ListGraphSnapshotsOutput{
			GraphSnapshots: []neptunegraphtypes.GraphSnapshotSummary{},
		}, nil)

	lister := &NeptuneGraphSnapshotLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testNeptuneGraphListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGraphSnapshot_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockNeptuneGraphClient)

	snap := &NeptuneGraphSnapshot{
		svc:        mockClient,
		SnapshotID: ptr.String("gs-abc123"),
	}

	mockClient.On("DeleteGraphSnapshot", mock.Anything, &neptunegraph.DeleteGraphSnapshotInput{
		SnapshotIdentifier: snap.SnapshotID,
	}).Return(&neptunegraph.DeleteGraphSnapshotOutput{}, nil)

	a.NoError(snap.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_NeptuneGraphSnapshot_Properties(t *testing.T) {
	a := assert.New(t)

	snap := NeptuneGraphSnapshot{
		SnapshotID: ptr.String("gs-abc123"),
		Name:       ptr.String("my-snapshot"),
		Status:     "AVAILABLE",
	}

	props := snap.Properties()
	a.Equal("gs-abc123", props.Get("SnapshotID"))
	a.Equal("my-snapshot", props.Get("Name"))
	a.Equal("AVAILABLE", props.Get("Status"))
}

func Test_Mock_NeptuneGraphSnapshot_String(t *testing.T) {
	a := assert.New(t)
	snap := NeptuneGraphSnapshot{SnapshotID: ptr.String("gs-abc123")}
	a.Equal("gs-abc123", snap.String())
}
