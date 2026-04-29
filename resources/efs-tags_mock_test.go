package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_EFSTags_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-123")},
			},
		}, nil)
	mockClient.On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(&efs.ListTagsForResourceOutput{
			Tags: []efstypes.Tag{
				{Key: ptr.String("Name"), Value: ptr.String("test")},
				{Key: ptr.String("aws:createdBy"), Value: ptr.String("system")},
			},
		}, nil)
	lister := &EFSTagsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*EFSTags)
	a.Equal("fs-123", r.String())
	a.Equal(1, *r.TagCount)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSTags_List_OnlyAwsTags(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-456")},
			},
		}, nil)
	mockClient.On("ListTagsForResource", mock.Anything, mock.Anything).
		Return(&efs.ListTagsForResourceOutput{
			Tags: []efstypes.Tag{
				{Key: ptr.String("aws:createdBy"), Value: ptr.String("system")},
			},
		}, nil)
	lister := &EFSTagsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSTags_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{},
		}, nil)
	lister := &EFSTagsLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSTags_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	r := &EFSTags{
		svc:          mockClient,
		FileSystemID: ptr.String("fs-123"),
		tagKeys:      []string{"Name", "Env"},
	}
	mockClient.On("UntagResource", mock.Anything, &efs.UntagResourceInput{
		ResourceId: r.FileSystemID,
		TagKeys:    []string{"Name", "Env"},
	}).Return(&efs.UntagResourceOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSTags_Properties(t *testing.T) {
	a := assert.New(t)
	tagCount := 2
	r := EFSTags{
		FileSystemID: ptr.String("fs-123"),
		TagCount:     &tagCount,
	}
	props := r.Properties()
	a.Equal("fs-123", props.Get("FileSystemID"))
	a.Equal("2", props.Get("TagCount"))
}

func Test_Mock_EFSTags_String(t *testing.T) {
	a := assert.New(t)
	r := &EFSTags{FileSystemID: ptr.String("fs-123")}
	a.Equal("fs-123", r.String())
}
