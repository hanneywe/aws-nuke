package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/resiliencehub"
	resiliencehubtypes "github.com/aws/aws-sdk-go-v2/service/resiliencehub/types"
)

func Test_Mock_ResilienceHubApp_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockResilienceHubClient)
	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&resiliencehub.ListAppsOutput{
			AppSummaries: []resiliencehubtypes.AppSummary{
				{AppArn: ptr.String("arn:aws:resiliencehub:us-east-1:123456789012:app/my-app"), Name: ptr.String("my-app")},
			},
		}, nil)
	lister := &ResilienceHubAppLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testResilienceHubListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	app := resources[0].(*ResilienceHubApp)
	a.Equal("my-app", *app.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ResilienceHubApp_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockResilienceHubClient)
	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&resiliencehub.ListAppsOutput{AppSummaries: []resiliencehubtypes.AppSummary{}}, nil)
	lister := &ResilienceHubAppLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testResilienceHubListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ResilienceHubApp_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockResilienceHubClient)
	app := &ResilienceHubApp{svc: mockClient, AppArn: ptr.String("arn:aws:resiliencehub:us-east-1:123456789012:app/my-app")}
	mockClient.On("DeleteApp", mock.Anything, &resiliencehub.DeleteAppInput{AppArn: app.AppArn}).
		Return(&resiliencehub.DeleteAppOutput{}, nil)
	a.NoError(app.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ResilienceHubApp_Properties(t *testing.T) {
	a := assert.New(t)
	app := ResilienceHubApp{AppArn: ptr.String("arn:aws:resiliencehub:us-east-1:123456789012:app/my-app"), Name: ptr.String("my-app")}
	a.Equal("my-app", app.Properties().Get("Name"))
}

func Test_Mock_ResilienceHubApp_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-app", (&ResilienceHubApp{Name: ptr.String("my-app")}).String())
}
