package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/synthetics"
	syntheticstypes "github.com/aws/aws-sdk-go-v2/service/synthetics/types"
)

func Test_Mock_SyntheticsGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSyntheticsClient)
	mockClient.On("ListGroups", mock.Anything, mock.Anything).
		Return(&synthetics.ListGroupsOutput{
			Groups: []syntheticstypes.GroupSummary{
				{Id: ptr.String("group-123"), Name: ptr.String("my-group")},
			},
		}, nil)
	lister := &SyntheticsGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSyntheticsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	g := resources[0].(*SyntheticsGroup)
	a.Equal("my-group", *g.Name)
	a.Equal("group-123", *g.GroupID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SyntheticsGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSyntheticsClient)
	mockClient.On("ListGroups", mock.Anything, mock.Anything).
		Return(&synthetics.ListGroupsOutput{Groups: []syntheticstypes.GroupSummary{}}, nil)
	lister := &SyntheticsGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSyntheticsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SyntheticsGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSyntheticsClient)
	g := &SyntheticsGroup{svc: mockClient, GroupID: ptr.String("group-123")}
	mockClient.On("DeleteGroup", mock.Anything, &synthetics.DeleteGroupInput{GroupIdentifier: g.GroupID}).
		Return(&synthetics.DeleteGroupOutput{}, nil)
	a.NoError(g.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SyntheticsGroup_Properties(t *testing.T) {
	a := assert.New(t)
	g := SyntheticsGroup{GroupID: ptr.String("group-123"), Name: ptr.String("my-group")}
	a.Equal("my-group", g.Properties().Get("Name"))
	a.Equal("group-123", g.Properties().Get("GroupId"))
}

func Test_Mock_SyntheticsGroup_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-group", (&SyntheticsGroup{Name: ptr.String("my-group")}).String())
}
