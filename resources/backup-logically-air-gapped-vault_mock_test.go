package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/backup"
	backuptypes "github.com/aws/aws-sdk-go-v2/service/backup/types"
)

func Test_Mock_BackupLogicallyAirGappedVault_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBackupClient)
	now := time.Now()
	mockClient.On("ListBackupVaults", mock.Anything, mock.Anything).
		Return(&backup.ListBackupVaultsOutput{
			BackupVaultList: []backuptypes.BackupVaultListMember{
				{
					BackupVaultName: ptr.String("my-vault"),
					BackupVaultArn:  ptr.String("arn:aws:backup:us-east-1:123456789012:backup-vault:my-vault"),
					CreationDate:    &now,
				},
			},
		}, nil)
	lister := &BackupLogicallyAirGappedVaultLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBackupListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	v := resources[0].(*BackupLogicallyAirGappedVault)
	a.Equal("my-vault", *v.Name)
	a.Equal("arn:aws:backup:us-east-1:123456789012:backup-vault:my-vault", *v.ARN)
	a.Equal(now, *v.CreationDate)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BackupLogicallyAirGappedVault_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBackupClient)
	mockClient.On("ListBackupVaults", mock.Anything, mock.Anything).
		Return(&backup.ListBackupVaultsOutput{
			BackupVaultList: []backuptypes.BackupVaultListMember{},
		}, nil)
	lister := &BackupLogicallyAirGappedVaultLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBackupListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BackupLogicallyAirGappedVault_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBackupClient)
	v := &BackupLogicallyAirGappedVault{
		svc:  mockClient,
		Name: ptr.String("my-vault"),
	}
	mockClient.On("DeleteBackupVault", mock.Anything, &backup.DeleteBackupVaultInput{
		BackupVaultName: v.Name,
	}).Return(&backup.DeleteBackupVaultOutput{}, nil)
	a.NoError(v.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_BackupLogicallyAirGappedVault_Properties(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	v := BackupLogicallyAirGappedVault{
		Name:         ptr.String("my-vault"),
		ARN:          ptr.String("arn:aws:backup:us-east-1:123456789012:backup-vault:my-vault"),
		CreationDate: &now,
	}
	props := v.Properties()
	a.Equal("my-vault", props.Get("Name"))
	a.Equal("arn:aws:backup:us-east-1:123456789012:backup-vault:my-vault", props.Get("ARN"))
}

func Test_Mock_BackupLogicallyAirGappedVault_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-vault", (&BackupLogicallyAirGappedVault{Name: ptr.String("my-vault")}).String())
}

func Test_Mock_BackupLogicallyAirGappedVault_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	v := BackupLogicallyAirGappedVault{Name: ptr.String("aws/my-managed-vault")}
	a.Error(v.Filter())
}

func Test_Mock_BackupLogicallyAirGappedVault_Filter_Normal(t *testing.T) {
	a := assert.New(t)
	v := BackupLogicallyAirGappedVault{Name: ptr.String("my-vault")}
	a.NoError(v.Filter())
}
