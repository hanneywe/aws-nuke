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

func Test_Mock_GrafanaWorkspaceAPIKey_List(t *testing.T) {
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
			WorkspaceId: ptr.String("g-abc123"),
		}, nil)

	mockClient.On("ListWorkspaceServiceAccountTokens", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspaceServiceAccountTokensOutput{
			ServiceAccountTokens: []grafanatypes.ServiceAccountTokenSummary{
				{Id: ptr.String("token-1"), Name: ptr.String("my-key")},
			},
		}, nil)

	lister := &GrafanaWorkspaceAPIKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*GrafanaWorkspaceAPIKey)
	a.Equal("g-abc123", *r.WorkspaceID)
	a.Equal("sa-1", *r.ServiceAccountID)
	a.Equal("token-1", *r.TokenID)
	a.Equal("my-key", *r.KeyName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceAPIKey_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	mockClient.On("ListWorkspaces", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspacesOutput{
			Workspaces: []grafanatypes.WorkspaceSummary{},
		}, nil)

	lister := &GrafanaWorkspaceAPIKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceAPIKey_List_MultipleWorkspaces(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	mockClient.On("ListWorkspaces", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspacesOutput{
			Workspaces: []grafanatypes.WorkspaceSummary{
				{Id: ptr.String("g-ws1")},
				{Id: ptr.String("g-ws2")},
			},
		}, nil)

	mockClient.On("ListWorkspaceServiceAccounts", mock.Anything, &grafana.ListWorkspaceServiceAccountsInput{
		WorkspaceId: ptr.String("g-ws1"),
	}).Return(&grafana.ListWorkspaceServiceAccountsOutput{
		ServiceAccounts: []grafanatypes.ServiceAccountSummary{
			{Id: ptr.String("sa-1"), Name: ptr.String("sa-one")},
		},
		WorkspaceId: ptr.String("g-ws1"),
	}, nil)

	mockClient.On("ListWorkspaceServiceAccounts", mock.Anything, &grafana.ListWorkspaceServiceAccountsInput{
		WorkspaceId: ptr.String("g-ws2"),
	}).Return(&grafana.ListWorkspaceServiceAccountsOutput{
		ServiceAccounts: []grafanatypes.ServiceAccountSummary{},
		WorkspaceId:     ptr.String("g-ws2"),
	}, nil)

	mockClient.On("ListWorkspaceServiceAccountTokens", mock.Anything, &grafana.ListWorkspaceServiceAccountTokensInput{
		WorkspaceId:      ptr.String("g-ws1"),
		ServiceAccountId: ptr.String("sa-1"),
	}).Return(&grafana.ListWorkspaceServiceAccountTokensOutput{
		ServiceAccountTokens: []grafanatypes.ServiceAccountTokenSummary{
			{Id: ptr.String("tok-1"), Name: ptr.String("key-1")},
		},
	}, nil)

	lister := &GrafanaWorkspaceAPIKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*GrafanaWorkspaceAPIKey)
	a.Equal("g-ws1", *r.WorkspaceID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceAPIKey_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	r := &GrafanaWorkspaceAPIKey{
		svc:              mockClient,
		WorkspaceID:      ptr.String("g-abc123"),
		ServiceAccountID: ptr.String("sa-1"),
		TokenID:          ptr.String("token-1"),
		KeyName:          ptr.String("my-key"),
	}

	mockClient.On("DeleteWorkspaceServiceAccountToken", mock.Anything,
		&grafana.DeleteWorkspaceServiceAccountTokenInput{
			WorkspaceId:      r.WorkspaceID,
			ServiceAccountId: r.ServiceAccountID,
			TokenId:          r.TokenID,
		}).Return(&grafana.DeleteWorkspaceServiceAccountTokenOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceAPIKey_Properties(t *testing.T) {
	a := assert.New(t)
	r := &GrafanaWorkspaceAPIKey{
		WorkspaceID:      ptr.String("g-abc123"),
		ServiceAccountID: ptr.String("sa-1"),
		TokenID:          ptr.String("token-1"),
		KeyName:          ptr.String("my-key"),
	}
	props := r.Properties()
	a.Equal("g-abc123", props.Get("WorkspaceId"))
	a.Equal("my-key", props.Get("KeyName"))
}

func Test_Mock_GrafanaWorkspaceAPIKey_String(t *testing.T) {
	a := assert.New(t)
	r := &GrafanaWorkspaceAPIKey{
		WorkspaceID: ptr.String("g-abc123"),
		KeyName:     ptr.String("my-key"),
	}
	a.Equal("g-abc123/my-key", r.String())
}
