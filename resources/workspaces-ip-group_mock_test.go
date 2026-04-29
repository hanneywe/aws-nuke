package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/workspaces"
	workspacestypes "github.com/aws/aws-sdk-go-v2/service/workspaces/types"
)

func Test_Mock_WorkSpacesIPGroup_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockWorkSpacesV2Client)

	mockClient.On("DescribeIpGroups", mock.Anything, mock.Anything).
		Return(&workspaces.DescribeIpGroupsOutput{
			Result: []workspacestypes.WorkspacesIpGroup{
				{
					GroupId:   ptr.String("wsipg-12345"),
					GroupName: ptr.String("my-ip-group"),
					GroupDesc: ptr.String("test group"),
				},
			},
		}, nil)

	lister := &WorkSpacesIPGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testWorkSpacesV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*WorkSpacesIPGroup)
	a.Equal("wsipg-12345", *r.GroupID)
	a.Equal("my-ip-group", *r.GroupName)
	a.Equal("test group", *r.GroupDesc)
	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesIPGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockWorkSpacesV2Client)

	mockClient.On("DescribeIpGroups", mock.Anything, mock.Anything).
		Return(&workspaces.DescribeIpGroupsOutput{
			Result: []workspacestypes.WorkspacesIpGroup{},
		}, nil)

	lister := &WorkSpacesIPGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testWorkSpacesV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesIPGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockWorkSpacesV2Client)

	r := &WorkSpacesIPGroup{
		svc:     mockClient,
		GroupID: ptr.String("wsipg-12345"),
	}

	mockClient.On("DeleteIpGroup", mock.Anything,
		&workspaces.DeleteIpGroupInput{
			GroupId: r.GroupID,
		}).Return(&workspaces.DeleteIpGroupOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_WorkSpacesIPGroup_Properties(t *testing.T) {
	a := assert.New(t)
	r := &WorkSpacesIPGroup{
		GroupID:   ptr.String("wsipg-12345"),
		GroupName: ptr.String("my-ip-group"),
		GroupDesc: ptr.String("test group"),
	}
	props := r.Properties()
	a.Equal("wsipg-12345", props.Get("GroupId"))
	a.Equal("my-ip-group", props.Get("GroupName"))
	a.Equal("test group", props.Get("GroupDesc"))
}

func Test_Mock_WorkSpacesIPGroup_String(t *testing.T) {
	a := assert.New(t)
	r := &WorkSpacesIPGroup{
		GroupID: ptr.String("wsipg-12345"),
	}
	a.Equal("wsipg-12345", r.String())
}
