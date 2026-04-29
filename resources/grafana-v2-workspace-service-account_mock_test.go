package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/grafana"
	grafanatypes "github.com/aws/aws-sdk-go-v2/service/grafana/types"
)

func Test_Mock_GrafanaWorkspaceServiceAccount_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	mockClient.On("ListWorkspaces", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspacesOutput{
			Workspaces: []grafanatypes.WorkspaceSummary{
				{Id: ptr.String("g-abc123")},
			},
		}, nil)

	mockClient.On("ListWorkspaceServiceAccounts", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspaceServiceAccountsOutput{
			ServiceAccounts: []grafanatypes.ServiceAccountSummary{
				{Id: ptr.String("sa-1"), Name: ptr.String("my-sa")},
			},
		}, nil)

	lister := &GrafanaWorkspaceServiceAccountLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*GrafanaWorkspaceServiceAccount)
	a.Equal("g-abc123", *r.WorkspaceID)
	a.Equal("sa-1", *r.ServiceAccountID)
	a.Equal("my-sa", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceServiceAccount_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	mockClient.On("ListWorkspaces", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspacesOutput{
			Workspaces: []grafanatypes.WorkspaceSummary{},
		}, nil)

	lister := &GrafanaWorkspaceServiceAccountLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceServiceAccount_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	r := &GrafanaWorkspaceServiceAccount{
		svc:              mockClient,
		WorkspaceID:      ptr.String("g-abc123"),
		ServiceAccountID: ptr.String("sa-1"),
		Name:             ptr.String("my-sa"),
	}

	mockClient.On("DeleteWorkspaceServiceAccount", mock.Anything,
		&grafana.DeleteWorkspaceServiceAccountInput{
			WorkspaceId:      r.WorkspaceID,
			ServiceAccountId: r.ServiceAccountID,
		}).Return(&grafana.DeleteWorkspaceServiceAccountOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceServiceAccount_Properties(t *testing.T) {
	a := assert.New(t)
	r := &GrafanaWorkspaceServiceAccount{
		WorkspaceID:      ptr.String("g-abc123"),
		ServiceAccountID: ptr.String("sa-1"),
		Name:             ptr.String("my-sa"),
	}
	props := r.Properties()
	a.Equal("g-abc123", props.Get("WorkspaceId"))
	a.Equal("sa-1", props.Get("ServiceAccountId"))
	a.Equal("my-sa", props.Get("Name"))
}

func Test_Mock_GrafanaWorkspaceServiceAccount_String(t *testing.T) {
	a := assert.New(t)
	r := &GrafanaWorkspaceServiceAccount{
		Name: ptr.String("my-sa"),
	}
	a.Equal("my-sa", r.String())
}
