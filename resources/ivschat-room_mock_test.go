package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivschat"
	ivschattypes "github.com/aws/aws-sdk-go-v2/service/ivschat/types"
)

func Test_Mock_IVSChatRoom_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSChatClient)

	mockClient.On("ListRooms", mock.Anything, mock.Anything).
		Return(&ivschat.ListRoomsOutput{
			Rooms: []ivschattypes.RoomSummary{
				{
					Arn:  ptr.String("arn:aws:ivschat:us-east-1:123456789012:room/abc123"),
					Name: ptr.String("my-room"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &IVSChatRoomLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSChatListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	room := resources[0].(*IVSChatRoom)
	a.Equal("my-room", *room.Name)
	a.Equal("test", room.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChatRoom_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSChatClient)

	mockClient.On("ListRooms", mock.Anything, mock.Anything).
		Return(&ivschat.ListRoomsOutput{
			Rooms: []ivschattypes.RoomSummary{},
		}, nil)

	lister := &IVSChatRoomLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSChatListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChatRoom_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSChatClient)

	room := &IVSChatRoom{
		svc: mockClient,
		ARN: ptr.String("arn:aws:ivschat:us-east-1:123456789012:room/abc123"),
	}

	mockClient.On("DeleteRoom", mock.Anything, &ivschat.DeleteRoomInput{
		Identifier: room.ARN,
	}).Return(&ivschat.DeleteRoomOutput{}, nil)

	a.NoError(room.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSChatRoom_Properties(t *testing.T) {
	a := assert.New(t)

	room := IVSChatRoom{
		ARN:  ptr.String("arn:aws:ivschat:us-east-1:123456789012:room/abc123"),
		Name: ptr.String("my-room"),
		Tags: map[string]string{"env": "test"},
	}

	props := room.Properties()
	a.Equal("my-room", props.Get("Name"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_IVSChatRoom_String(t *testing.T) {
	a := assert.New(t)
	room := IVSChatRoom{Name: ptr.String("my-room")}
	a.Equal("my-room", room.String())
}
