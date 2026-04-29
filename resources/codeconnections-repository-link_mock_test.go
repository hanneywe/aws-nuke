package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codeconnections"
	codeconnectionstypes "github.com/aws/aws-sdk-go-v2/service/codeconnections/types"
)

func Test_Mock_CodeConnectionsRepositoryLink_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeConnectionsClient)
	mockClient.On("ListRepositoryLinks", mock.Anything, mock.Anything).
		Return(&codeconnections.ListRepositoryLinksOutput{
			RepositoryLinks: []codeconnectionstypes.RepositoryLinkInfo{
				{
					RepositoryLinkId: ptr.String("repo-link-123"),
					RepositoryName:   ptr.String("my-repo"),
					ProviderType:     codeconnectionstypes.ProviderTypeGithub,
					ConnectionArn:    ptr.String("arn:aws:codeconnections:us-east-1:123456789012:connection/abc-123"),
				},
			},
		}, nil)
	lister := &CodeConnectionsRepositoryLinkLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeConnectionsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*CodeConnectionsRepositoryLink)
	a.Equal("repo-link-123", *r.RepositoryLinkID)
	a.Equal("my-repo", *r.RepositoryName)
	a.Equal(codeconnectionstypes.ProviderTypeGithub, r.ProviderType)
	a.Equal("arn:aws:codeconnections:us-east-1:123456789012:connection/abc-123", *r.ConnectionArn)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeConnectionsRepositoryLink_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeConnectionsClient)
	mockClient.On("ListRepositoryLinks", mock.Anything, mock.Anything).
		Return(&codeconnections.ListRepositoryLinksOutput{
			RepositoryLinks: []codeconnectionstypes.RepositoryLinkInfo{},
		}, nil)
	lister := &CodeConnectionsRepositoryLinkLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeConnectionsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeConnectionsRepositoryLink_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeConnectionsClient)
	r := &CodeConnectionsRepositoryLink{
		svc:              mockClient,
		RepositoryLinkID: ptr.String("repo-link-123"),
	}
	mockClient.On("DeleteRepositoryLink", mock.Anything, &codeconnections.DeleteRepositoryLinkInput{
		RepositoryLinkId: r.RepositoryLinkID,
	}).Return(&codeconnections.DeleteRepositoryLinkOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeConnectionsRepositoryLink_Properties(t *testing.T) {
	a := assert.New(t)
	r := CodeConnectionsRepositoryLink{
		RepositoryLinkID: ptr.String("repo-link-123"),
		RepositoryName:   ptr.String("my-repo"),
		ProviderType:     codeconnectionstypes.ProviderTypeGithub,
		ConnectionArn:    ptr.String("arn:aws:codeconnections:us-east-1:123456789012:connection/abc-123"),
	}
	props := r.Properties()
	a.Equal("repo-link-123", props.Get("RepositoryLinkID"))
	a.Equal("my-repo", props.Get("RepositoryName"))
	a.Equal("arn:aws:codeconnections:us-east-1:123456789012:connection/abc-123", props.Get("ConnectionArn"))
}

func Test_Mock_CodeConnectionsRepositoryLink_String(t *testing.T) {
	a := assert.New(t)
	r := &CodeConnectionsRepositoryLink{RepositoryLinkID: ptr.String("repo-link-123")}
	a.Equal("repo-link-123", r.String())
}
