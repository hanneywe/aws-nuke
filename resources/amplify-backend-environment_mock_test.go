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

func Test_Mock_AmplifyBackendEnvironment_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)

	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{
			Apps: []amplifytypes.App{
				{AppId: ptr.String("app-123")},
			},
		}, nil)

	mockClient.On("ListBackendEnvironments", mock.Anything, mock.Anything).
		Return(&amplify.ListBackendEnvironmentsOutput{
			BackendEnvironments: []amplifytypes.BackendEnvironment{
				{EnvironmentName: ptr.String("staging")},
			},
		}, nil)

	lister := &AmplifyBackendEnvironmentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	env := resources[0].(*AmplifyBackendEnvironment)
	a.Equal("staging", *env.EnvironmentName)
	a.Equal("app-123", *env.AppID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyBackendEnvironment_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{Apps: []amplifytypes.App{}}, nil)
	lister := &AmplifyBackendEnvironmentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyBackendEnvironment_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	env := &AmplifyBackendEnvironment{svc: mockClient, AppID: ptr.String("app-123"), EnvironmentName: ptr.String("staging")}
	mockClient.On("DeleteBackendEnvironment", mock.Anything, &amplify.DeleteBackendEnvironmentInput{
		AppId: env.AppID, EnvironmentName: env.EnvironmentName,
	}).Return(&amplify.DeleteBackendEnvironmentOutput{}, nil)
	a.NoError(env.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyBackendEnvironment_Properties(t *testing.T) {
	a := assert.New(t)
	env := AmplifyBackendEnvironment{EnvironmentName: ptr.String("staging"), AppID: ptr.String("app-123")}
	a.Equal("staging", env.Properties().Get("EnvironmentName"))
	a.Equal("app-123", env.Properties().Get("AppId"))
}

func Test_Mock_AmplifyBackendEnvironment_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("staging", (&AmplifyBackendEnvironment{EnvironmentName: ptr.String("staging")}).String())
}
