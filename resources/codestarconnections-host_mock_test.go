package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/codestarconnections"
	codestarconnectionstypes "github.com/aws/aws-sdk-go-v2/service/codestarconnections/types"
)

func Test_Mock_CodeStarConnectionsHost_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeStarConnectionsClient)
	mockClient.On("ListHosts", mock.Anything, mock.Anything).
		Return(&codestarconnections.ListHostsOutput{
			Hosts: []codestarconnectionstypes.Host{
				{HostArn: ptr.String("arn:aws:codestar-connections:us-east-1:123456789012:host/my-host"), Name: ptr.String("my-host")},
			},
		}, nil)
	lister := &CodeStarConnectionsHostLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeStarConnectionsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-host", *resources[0].(*CodeStarConnectionsHost).Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeStarConnectionsHost_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeStarConnectionsClient)
	mockClient.On("ListHosts", mock.Anything, mock.Anything).
		Return(&codestarconnections.ListHostsOutput{Hosts: []codestarconnectionstypes.Host{}}, nil)
	lister := &CodeStarConnectionsHostLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCodeStarConnectionsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeStarConnectionsHost_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCodeStarConnectionsClient)
	r := &CodeStarConnectionsHost{
		svc:     mockClient,
		HostArn: ptr.String("arn:aws:codestar-connections:us-east-1:123456789012:host/my-host"),
		Name:    ptr.String("my-host"),
	}
	mockClient.On("DeleteHost", mock.Anything, &codestarconnections.DeleteHostInput{HostArn: r.HostArn}).
		Return(&codestarconnections.DeleteHostOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CodeStarConnectionsHost_Properties(t *testing.T) {
	a := assert.New(t)
	r := CodeStarConnectionsHost{HostArn: ptr.String("arn"), Name: ptr.String("my-host")}
	a.Equal("my-host", r.Properties().Get("Name"))
	a.Equal("arn", r.Properties().Get("HostArn"))
}

func Test_Mock_CodeStarConnectionsHost_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-host", (&CodeStarConnectionsHost{Name: ptr.String("my-host")}).String())
}
