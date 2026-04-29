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

func Test_Mock_IoTMitigationAction_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListMitigationActions", mock.Anything, mock.Anything).
		Return(&iot.ListMitigationActionsOutput{
			ActionIdentifiers: []iottypes.MitigationActionIdentifier{
				{
					ActionName: ptr.String("my-mitigation-action"),
					ActionArn:  ptr.String("arn:aws:iot:us-east-1:123456789012:mitigationaction/my-mitigation-action"),
				},
			},
		}, nil)

	lister := &IoTMitigationActionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	action := resources[0].(*IoTMitigationAction)
	assertions.Equal("my-mitigation-action", *action.ActionName)
	assertions.Equal("arn:aws:iot:us-east-1:123456789012:mitigationaction/my-mitigation-action", *action.ActionArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTMitigationAction_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListMitigationActions", mock.Anything, mock.Anything).
		Return(&iot.ListMitigationActionsOutput{
			ActionIdentifiers: []iottypes.MitigationActionIdentifier{},
		}, nil)

	lister := &IoTMitigationActionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTMitigationAction_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTClient)

	action := &IoTMitigationAction{
		svc:        mockClient,
		ActionName: ptr.String("my-mitigation-action"),
	}

	mockClient.On("DeleteMitigationAction", mock.Anything, &iot.DeleteMitigationActionInput{
		ActionName: action.ActionName,
	}).Return(&iot.DeleteMitigationActionOutput{}, nil)

	assertions.NoError(action.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTMitigationAction_Properties(t *testing.T) {
	assertions := assert.New(t)

	action := IoTMitigationAction{
		ActionName: ptr.String("my-mitigation-action"),
		ActionArn:  ptr.String("arn:aws:iot:us-east-1:123456789012:mitigationaction/my-mitigation-action"),
	}

	props := action.Properties()
	assertions.Equal("my-mitigation-action", props.Get("ActionName"))
	assertions.Equal("arn:aws:iot:us-east-1:123456789012:mitigationaction/my-mitigation-action", props.Get("ActionArn"))
}

func Test_Mock_IoTMitigationAction_String(t *testing.T) {
	assertions := assert.New(t)
	action := IoTMitigationAction{ActionName: ptr.String("my-mitigation-action")}
	assertions.Equal("my-mitigation-action", action.String())
}
