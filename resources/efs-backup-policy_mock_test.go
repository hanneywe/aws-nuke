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

func Test_Mock_EFSBackupPolicy_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-12345")},
				{FileSystemId: ptr.String("fs-67890")},
			},
		}, nil)

	mockClient.On("DescribeBackupPolicy", mock.Anything, &efs.DescribeBackupPolicyInput{
		FileSystemId: ptr.String("fs-12345"),
	}).Return(&efs.DescribeBackupPolicyOutput{
		BackupPolicy: &efstypes.BackupPolicy{
			Status: efstypes.StatusEnabled,
		},
	}, nil)

	mockClient.On("DescribeBackupPolicy", mock.Anything, &efs.DescribeBackupPolicyInput{
		FileSystemId: ptr.String("fs-67890"),
	}).Return(&efs.DescribeBackupPolicyOutput{
		BackupPolicy: &efstypes.BackupPolicy{
			Status: efstypes.StatusEnabled,
		},
	}, nil)

	lister := &EFSBackupPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 2)

	bp := resources[0].(*EFSBackupPolicy)
	a.Equal("fs-12345", *bp.FileSystemID)
	a.Equal("ENABLED", *bp.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSBackupPolicy_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{},
		}, nil)

	lister := &EFSBackupPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSBackupPolicy_List_BackupDisabled(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-12345")},
			},
		}, nil)

	mockClient.On("DescribeBackupPolicy", mock.Anything, &efs.DescribeBackupPolicyInput{
		FileSystemId: ptr.String("fs-12345"),
	}).Return(&efs.DescribeBackupPolicyOutput{
		BackupPolicy: &efstypes.BackupPolicy{
			Status: efstypes.StatusDisabled,
		},
	}, nil)

	lister := &EFSBackupPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSBackupPolicy_List_PolicyNotFound(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	mockClient.On("DescribeFileSystems", mock.Anything, mock.Anything).
		Return(&efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{FileSystemId: ptr.String("fs-12345")},
			},
		}, nil)

	mockClient.On("DescribeBackupPolicy", mock.Anything, &efs.DescribeBackupPolicyInput{
		FileSystemId: ptr.String("fs-12345"),
	}).Return((*efs.DescribeBackupPolicyOutput)(nil), &efstypes.PolicyNotFound{
		Message: ptr.String("no backup policy"),
	})

	lister := &EFSBackupPolicyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEFSV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSBackupPolicy_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEFSV2Client)

	bp := &EFSBackupPolicy{
		svc:          mockClient,
		FileSystemID: ptr.String("fs-12345"),
		Status:       ptr.String("ENABLED"),
	}

	mockClient.On("PutBackupPolicy", mock.Anything, &efs.PutBackupPolicyInput{
		FileSystemId: ptr.String("fs-12345"),
		BackupPolicy: &efstypes.BackupPolicy{
			Status: efstypes.StatusDisabled,
		},
	}).Return(&efs.PutBackupPolicyOutput{}, nil)

	err := bp.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EFSBackupPolicy_Properties(t *testing.T) {
	a := assert.New(t)

	bp := EFSBackupPolicy{
		FileSystemID: ptr.String("fs-12345"),
		Status:       ptr.String("ENABLED"),
	}

	props := bp.Properties()
	a.Equal("fs-12345", props.Get("FileSystemId"))
	a.Equal("ENABLED", props.Get("Status"))
}

func Test_Mock_EFSBackupPolicy_String(t *testing.T) {
	a := assert.New(t)
	bp := EFSBackupPolicy{FileSystemID: ptr.String("fs-12345")}
	a.Equal("fs-12345 -> BackupPolicy", bp.String())
}
