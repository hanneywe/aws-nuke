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

func Test_Mock_EFSLifecycleConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-123")},
			},
		}, nil)
	mockClient.On("DescribeLifecycleConfiguration", mock.Anything, mock.Anything).
		Return(&efs.DescribeLifecycleConfigurationOutput{
			LifecyclePolicies: []efstypes.LifecyclePolicy{
				{TransitionToIA: efstypes.TransitionToIARulesAfter30Days},
			},
		}, nil)
	lister := &EFSLifecycleConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("fs-123", resources[0].(*EFSLifecycleConfiguration).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSLifecycleConfiguration_List_EmptyConfig(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-456")},
			},
		}, nil)
	mockClient.On("DescribeLifecycleConfiguration", mock.Anything, mock.Anything).
		Return(&efs.DescribeLifecycleConfigurationOutput{
			LifecyclePolicies: []efstypes.LifecyclePolicy{},
		}, nil)
	lister := &EFSLifecycleConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSLifecycleConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)
	r := &EFSLifecycleConfiguration{
		svc:          mockClient,
		FileSystemID: ptr.String("fs-123"),
	}
	mockClient.On("PutLifecycleConfiguration", mock.Anything, &efs.PutLifecycleConfigurationInput{
		FileSystemId:      r.FileSystemID,
		LifecyclePolicies: []efstypes.LifecyclePolicy{},
	}).Return(&efs.PutLifecycleConfigurationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSLifecycleConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := EFSLifecycleConfiguration{
		FileSystemID: ptr.String("fs-123"),
	}
	props := r.Properties()
	a.Equal("fs-123", props.Get("FileSystemID"))
}

func Test_Mock_EFSLifecycleConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &EFSLifecycleConfiguration{FileSystemID: ptr.String("fs-123")}
	a.Equal("fs-123", r.String())
}
