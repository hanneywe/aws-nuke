package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/backup"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockBackupClient struct {
	mock.Mock
}

func (m *mockBackupClient) ListBackupVaults(
	ctx context.Context, params *backup.ListBackupVaultsInput,
	_ ...func(*backup.Options),
) (*backup.ListBackupVaultsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*backup.ListBackupVaultsOutput), args.Error(1)
}

func (m *mockBackupClient) DeleteBackupVault(
	ctx context.Context, params *backup.DeleteBackupVaultInput,
	_ ...func(*backup.Options),
) (*backup.DeleteBackupVaultOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*backup.DeleteBackupVaultOutput), args.Error(1)
}

var testBackupListerOpts = &nuke.ListerOpts{}
