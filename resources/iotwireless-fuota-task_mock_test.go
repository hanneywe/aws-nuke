package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	iotwirelesstypes "github.com/aws/aws-sdk-go-v2/service/iotwireless/types"
)

func Test_Mock_IoTWirelessFuotaTask_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListFuotaTasks", mock.Anything, mock.Anything).
		Return(&iotwireless.ListFuotaTasksOutput{
			FuotaTaskList: []iotwirelesstypes.FuotaTask{
				{
					Id:   ptr.String("ft-11111"),
					Name: ptr.String("my-fuota-task"),
				},
				{
					Id:   ptr.String("ft-22222"),
					Name: ptr.String("another-fuota-task"),
				},
			},
		}, nil)

	lister := &IoTWirelessFuotaTaskLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	fuotaTask := resources[0].(*IoTWirelessFuotaTask)
	assertions.Equal("ft-11111", *fuotaTask.ID)
	assertions.Equal("my-fuota-task", *fuotaTask.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessFuotaTask_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListFuotaTasks", mock.Anything, mock.Anything).
		Return(&iotwireless.ListFuotaTasksOutput{
			FuotaTaskList: []iotwirelesstypes.FuotaTask{},
		}, nil)

	lister := &IoTWirelessFuotaTaskLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessFuotaTask_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	fuotaTask := &IoTWirelessFuotaTask{
		svc: mockClient,
		ID:  ptr.String("ft-11111"),
	}

	mockClient.On("DeleteFuotaTask", mock.Anything, &iotwireless.DeleteFuotaTaskInput{
		Id: fuotaTask.ID,
	}).Return(&iotwireless.DeleteFuotaTaskOutput{}, nil)

	assertions.NoError(fuotaTask.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessFuotaTask_Properties(t *testing.T) {
	assertions := assert.New(t)

	fuotaTask := IoTWirelessFuotaTask{
		ID:   ptr.String("ft-11111"),
		Name: ptr.String("my-fuota-task"),
	}

	props := fuotaTask.Properties()
	assertions.Equal("ft-11111", props.Get("Id"))
	assertions.Equal("my-fuota-task", props.Get("Name"))
}

func Test_Mock_IoTWirelessFuotaTask_String(t *testing.T) {
	assertions := assert.New(t)
	fuotaTask := IoTWirelessFuotaTask{ID: ptr.String("ft-11111")}
	assertions.Equal("ft-11111", fuotaTask.String())
}
