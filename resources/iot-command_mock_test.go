package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

func Test_Mock_IoTCommand_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListCommands", mock.Anything, mock.Anything).
		Return(&iot.ListCommandsOutput{
			Commands: []iottypes.CommandSummary{
				{
					CommandId:       ptr.String("cmd-12345"),
					CommandArn:      ptr.String("arn:aws:iot:us-east-1:123456789012:command/cmd-12345"),
					PendingDeletion: ptr.Bool(false),
				},
			},
		}, nil)

	lister := &IoTCommandLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cmd := resources[0].(*IoTCommand)
	a.Equal("cmd-12345", *cmd.CommandID)
	a.Equal("arn:aws:iot:us-east-1:123456789012:command/cmd-12345", *cmd.CommandArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTCommand_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListCommands", mock.Anything, mock.Anything).
		Return(&iot.ListCommandsOutput{
			Commands: []iottypes.CommandSummary{},
		}, nil)

	lister := &IoTCommandLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTCommand_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTClient)

	cmd := &IoTCommand{
		svc:       mockClient,
		CommandID: ptr.String("cmd-12345"),
	}

	mockClient.On("DeleteCommand", mock.Anything, &iot.DeleteCommandInput{
		CommandId: cmd.CommandID,
	}).Return(&iot.DeleteCommandOutput{}, nil)

	a.NoError(cmd.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTCommand_Filter_PendingDeletion(t *testing.T) {
	a := assert.New(t)

	cmd := IoTCommand{
		CommandID:       ptr.String("cmd-12345"),
		PendingDeletion: ptr.Bool(true),
	}
	a.Error(cmd.Filter())
	a.Contains(cmd.Filter().Error(), "already pending deletion")
}

func Test_Mock_IoTCommand_Filter_NotPending(t *testing.T) {
	a := assert.New(t)

	cmd := IoTCommand{
		CommandID:       ptr.String("cmd-12345"),
		PendingDeletion: ptr.Bool(false),
	}
	a.NoError(cmd.Filter())
}

func Test_Mock_IoTCommand_Filter_Nil(t *testing.T) {
	a := assert.New(t)

	cmd := IoTCommand{
		CommandID: ptr.String("cmd-12345"),
	}
	a.NoError(cmd.Filter())
}

func Test_Mock_IoTCommand_Properties(t *testing.T) {
	a := assert.New(t)

	cmd := IoTCommand{
		CommandID:  ptr.String("cmd-12345"),
		CommandArn: ptr.String("arn:aws:iot:us-east-1:123456789012:command/cmd-12345"),
	}

	props := cmd.Properties()
	a.Equal("cmd-12345", props.Get("CommandId"))
	a.Equal("arn:aws:iot:us-east-1:123456789012:command/cmd-12345", props.Get("CommandArn"))
}

func Test_Mock_IoTCommand_String(t *testing.T) {
	a := assert.New(t)
	cmd := IoTCommand{CommandID: ptr.String("cmd-12345")}
	a.Equal("cmd-12345", cmd.String())
}
