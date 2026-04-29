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

func Test_Mock_AmplifyBranch_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)

	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{
			Apps: []amplifytypes.App{
				{AppId: ptr.String("app-123")},
			},
		}, nil)

	mockClient.On("ListBranches", mock.Anything, mock.Anything).
		Return(&amplify.ListBranchesOutput{
			Branches: []amplifytypes.Branch{
				{BranchName: ptr.String("main"), BranchArn: ptr.String("arn:aws:amplify:us-east-1:123456789012:apps/app-123/branches/main")},
			},
		}, nil)

	lister := &AmplifyBranchLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	b := resources[0].(*AmplifyBranch)
	a.Equal("main", *b.BranchName)
	a.Equal("app-123", *b.AppID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyBranch_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	mockClient.On("ListApps", mock.Anything, mock.Anything).
		Return(&amplify.ListAppsOutput{Apps: []amplifytypes.App{}}, nil)
	lister := &AmplifyBranchLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAmplifyListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyBranch_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAmplifyClient)
	b := &AmplifyBranch{svc: mockClient, AppID: ptr.String("app-123"), BranchName: ptr.String("main")}
	mockClient.On("DeleteBranch", mock.Anything, &amplify.DeleteBranchInput{
		AppId: b.AppID, BranchName: b.BranchName,
	}).Return(&amplify.DeleteBranchOutput{}, nil)
	a.NoError(b.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AmplifyBranch_Properties(t *testing.T) {
	a := assert.New(t)
	b := AmplifyBranch{BranchName: ptr.String("main"), AppID: ptr.String("app-123")}
	a.Equal("main", b.Properties().Get("BranchName"))
	a.Equal("app-123", b.Properties().Get("AppId"))
}

func Test_Mock_AmplifyBranch_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("main", (&AmplifyBranch{BranchName: ptr.String("main")}).String())
}
