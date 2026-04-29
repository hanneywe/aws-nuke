package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
)

func Test_Mock_EFSReplicationConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeReplicationConfigurations", mock.Anything, mock.Anything).
		Return(&efs.DescribeReplicationConfigurationsOutput{
			Replications: []efstypes.ReplicationConfigurationDescription{
				{
					SourceFileSystemId:          ptr.String("fs-12345"),
					OriginalSourceFileSystemArn: ptr.String("arn:aws:elasticfilesystem:us-east-1:123456789012:file-system/fs-12345"),
					SourceFileSystemRegion:      ptr.String("us-east-1"),
				},
			},
		}, nil)

	lister := &EFSReplicationConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*EFSReplicationConfiguration)
	a.Equal("fs-12345", *r.FileSystemID)
	a.Equal("us-east-1", *r.SourceFileSystemRegion)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSReplicationConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeReplicationConfigurations", mock.Anything, mock.Anything).
		Return(&efs.DescribeReplicationConfigurationsOutput{
			Replications: []efstypes.ReplicationConfigurationDescription{},
		}, nil)

	lister := &EFSReplicationConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSReplicationConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	r := &EFSReplicationConfiguration{
		svc:          mockClient,
		FileSystemID: ptr.String("fs-12345"),
	}

	mockClient.On("DeleteReplicationConfiguration", mock.Anything, &efs.DeleteReplicationConfigurationInput{
		SourceFileSystemId: r.FileSystemID,
	}).Return(&efs.DeleteReplicationConfigurationOutput{}, nil)

	err := r.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSReplicationConfiguration_Properties(t *testing.T) {
	a := assert.New(t)

	r := EFSReplicationConfiguration{
		FileSystemID:                ptr.String("fs-12345"),
		OriginalSourceFileSystemArn: ptr.String("arn:aws:elasticfilesystem:us-east-1:123456789012:file-system/fs-12345"),
		SourceFileSystemRegion:      ptr.String("us-east-1"),
	}

	props := r.Properties()
	a.Equal("fs-12345", props.Get("FileSystemID"))
	a.Equal("arn:aws:elasticfilesystem:us-east-1:123456789012:file-system/fs-12345", props.Get("OriginalSourceFileSystemArn"))
	a.Equal("us-east-1", props.Get("SourceFileSystemRegion"))
}

func Test_Mock_EFSReplicationConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := EFSReplicationConfiguration{FileSystemID: ptr.String("fs-12345")}
	a.Equal("fs-12345", r.String())
}
