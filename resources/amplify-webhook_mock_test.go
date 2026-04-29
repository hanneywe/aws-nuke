package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/amplify"
	amplifytypes "github.com/aws/aws-sdk-go-v2/service/amplify/types"
)

func Test_Mock_AmplifyWebhook_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)

	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{
			Apps: []amplifytypes.App{
				{AppId: ptr.String("app-123")},
			},
		}, nil)

	mockClient.On("ListWebhooks", mock.Anything, mock.Anything).
		Return(&amplify.ListWebhooksOutput{
			Webhooks: []amplifytypes.Webhook{
				{WebhookId: ptr.String("wh-123"), WebhookArn: ptr.String("arn:aws:amplify:us-east-1:123456789012:apps/app-123/webhooks/wh-123")},
			},
		}, nil)

	lister := &AmplifyWebhookLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	wh := resources[0].(*AmplifyWebhook)
	a.Equal("wh-123", *wh.WebhookID)
	a.Equal("app-123", *wh.AppID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyWebhook_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{Apps: []amplifytypes.App{}}, nil)
	lister := &AmplifyWebhookLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyWebhook_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	wh := &AmplifyWebhook{svc: mockClient, WebhookID: ptr.String("wh-123")}
	mockClient.On("DeleteWebhook", mock.Anything, &amplify.DeleteWebhookInput{WebhookId: wh.WebhookID}).
		Return(&amplify.DeleteWebhookOutput{}, nil)
	a.NoError(wh.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyWebhook_Properties(t *testing.T) {
	a := assert.New(t)
	wh := AmplifyWebhook{WebhookID: ptr.String("wh-123"), AppID: ptr.String("app-123")}
	a.Equal("wh-123", wh.Properties().Get("WebhookId"))
	a.Equal("app-123", wh.Properties().Get("AppId"))
}

func Test_Mock_AmplifyWebhook_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("wh-123", (&AmplifyWebhook{WebhookID: ptr.String("wh-123")}).String())
}
