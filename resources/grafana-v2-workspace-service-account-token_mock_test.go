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

func Test_Mock_GrafanaWorkspaceServiceAccountToken_List(t *testing.T) {
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

	mockClient.On("ListWorkspaceServiceAccountTokens", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspaceServiceAccountTokensOutput{
			ServiceAccountTokens: []grafanatypes.ServiceAccountTokenSummary{
				{Id: ptr.String("token-1"), Name: ptr.String("my-token")},
			},
		}, nil)

	lister := &GrafanaWorkspaceServiceAccountTokenLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*GrafanaWorkspaceServiceAccountToken)
	a.Equal("g-abc123", *r.WorkspaceID)
	a.Equal("sa-1", *r.ServiceAccountID)
	a.Equal("token-1", *r.TokenID)
	a.Equal("my-token", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceServiceAccountToken_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	mockClient.On("ListWorkspaces", mock.Anything, mock.Anything).
		Return(&grafana.ListWorkspacesOutput{
			Workspaces: []grafanatypes.WorkspaceSummary{},
		}, nil)

	lister := &GrafanaWorkspaceServiceAccountTokenLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGrafanaV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GrafanaWorkspaceServiceAccountToken_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGrafanaV2Client)

	r := &GrafanaWorkspaceServiceAccountToken{
		svc:              mockClient,
		WorkspaceID:      ptr.String("g-abc123"),
		ServiceAccountID: ptr.String("sa-1"),
		TokenID:          ptr.String("token-1"),
		Name:             ptr.String("my-token"),
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

func Test_Mock_GrafanaWorkspaceServiceAccountToken_Properties(t *testing.T) {
	a := assert.New(t)
	r := &GrafanaWorkspaceServiceAccountToken{
		WorkspaceID:      ptr.String("g-abc123"),
		ServiceAccountID: ptr.String("sa-1"),
		TokenID:          ptr.String("token-1"),
		Name:             ptr.String("my-token"),
	}
	props := r.Properties()
	a.Equal("g-abc123", props.Get("WorkspaceId"))
	a.Equal("sa-1", props.Get("ServiceAccountId"))
	a.Equal("token-1", props.Get("TokenId"))
	a.Equal("my-token", props.Get("Name"))
}

func Test_Mock_GrafanaWorkspaceServiceAccountToken_String(t *testing.T) {
	a := assert.New(t)
	r := &GrafanaWorkspaceServiceAccountToken{
		Name: ptr.String("my-token"),
	}
	a.Equal("my-token", r.String())
}
