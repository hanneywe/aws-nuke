package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/chatbot"
	chatbottypes "github.com/aws/aws-sdk-go-v2/service/chatbot/types"
)

func Test_Mock_ChatbotCustomAction_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChatbotClient)

	arn := "arn:aws:chatbot::123456789012:custom-action/my-action"

	mockClient.
		On("ListCustomActions", mock.Anything, mock.Anything).
		Return(&chatbot.ListCustomActionsOutput{
			CustomActions: []string{arn},
		}, nil)

	mockClient.
		On("GetCustomAction", mock.Anything, mock.Anything).
		Return(&chatbot.GetCustomActionOutput{
			CustomAction: &chatbottypes.CustomAction{
				CustomActionArn: ptr.String(arn),
				ActionName:      ptr.String("my-action"),
			},
		}, nil)

	lister := &ChatbotCustomActionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testChatbotListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	action := resources[0].(*ChatbotCustomAction)
	a.Equal("my-action", *action.ActionName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ChatbotCustomAction_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChatbotClient)

	mockClient.
		On("ListCustomActions", mock.Anything, mock.Anything).
		Return(&chatbot.ListCustomActionsOutput{
			CustomActions: []string{},
		}, nil)

	lister := &ChatbotCustomActionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testChatbotListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ChatbotCustomAction_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockChatbotClient)

	action := &ChatbotCustomAction{
		svc:             mockClient,
		CustomActionArn: ptr.String("arn:aws:chatbot::123456789012:custom-action/my-action"),
	}

	mockClient.
		On("DeleteCustomAction", mock.Anything, &chatbot.DeleteCustomActionInput{
			CustomActionArn: action.CustomActionArn,
		}).
		Return(&chatbot.DeleteCustomActionOutput{}, nil)

	err := action.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ChatbotCustomAction_Properties(t *testing.T) {
	a := assert.New(t)

	action := ChatbotCustomAction{
		CustomActionArn: ptr.String("arn:aws:chatbot::123456789012:custom-action/my-action"),
		ActionName:      ptr.String("my-action"),
	}

	props := action.Properties()
	a.Equal("arn:aws:chatbot::123456789012:custom-action/my-action", props.Get("CustomActionArn"))
	a.Equal("my-action", props.Get("ActionName"))
}

func Test_Mock_ChatbotCustomAction_String(t *testing.T) {
	a := assert.New(t)

	action := ChatbotCustomAction{
		CustomActionArn: ptr.String("arn:aws:chatbot::123456789012:custom-action/my-action"),
	}

	a.Equal("arn:aws:chatbot::123456789012:custom-action/my-action", action.String())
}
