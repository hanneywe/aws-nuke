package resources

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/backup"
	backuptypes "github.com/aws/aws-sdk-go-v2/service/backup/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BackupLogicallyAirGappedVaultResource = "BackupLogicallyAirGappedVault"

func init() {
	registry.Register(&registry.Registration{
		Name:     BackupLogicallyAirGappedVaultResource,
		Scope:    nuke.Account,
		Resource: &BackupLogicallyAirGappedVault{},
		Lister:   &BackupLogicallyAirGappedVaultLister{},
	})
}

type BackupLogicallyAirGappedVaultLister struct {
	svc BackupClient
}

func (l *BackupLogicallyAirGappedVaultLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = backup.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := backup.NewListBackupVaultsPaginator(svc, &backup.ListBackupVaultsInput{
		ByVaultType: backuptypes.VaultTypeLogicallyAirGappedBackupVault,
	})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range resp.BackupVaultList {
			v := &resp.BackupVaultList[i]
			resources = append(resources, &BackupLogicallyAirGappedVault{
				svc:          svc,
				Name:         v.BackupVaultName,
				ARN:          v.BackupVaultArn,
				CreationDate: v.CreationDate,
			})
		}
	}

	return resources, nil
}

type BackupLogicallyAirGappedVault struct {
	svc          BackupClient
	Name         *string
	ARN          *string
	CreationDate *time.Time
}

func (r *BackupLogicallyAirGappedVault) Filter() error {
	if r.Name != nil && strings.HasPrefix(*r.Name, "aws/") {
		return fmt.Errorf("cannot delete AWS-managed vault")
	}
	return nil
}

func (r *BackupLogicallyAirGappedVault) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBackupVault(ctx, &backup.DeleteBackupVaultInput{
		BackupVaultName: r.Name,
	})
	return err
}

func (r *BackupLogicallyAirGappedVault) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BackupLogicallyAirGappedVault) String() string {
	return *r.Name
}
