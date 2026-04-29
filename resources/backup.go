package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/backup"
)

// BackupClient is the interface for the Backup SDK v2 client methods.
type BackupClient interface {
	ListBackupVaults(ctx context.Context, params *backup.ListBackupVaultsInput,
		optFns ...func(*backup.Options)) (*backup.ListBackupVaultsOutput, error)
	DeleteBackupVault(ctx context.Context, params *backup.DeleteBackupVaultInput,
		optFns ...func(*backup.Options)) (*backup.DeleteBackupVaultOutput, error)
}
