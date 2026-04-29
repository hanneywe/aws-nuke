package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EFSBackupPolicyResource = "EFSBackupPolicy"

func init() {
	registry.Register(&registry.Registration{
		Name:     EFSBackupPolicyResource,
		Scope:    nuke.Account,
		Resource: &EFSBackupPolicy{},
		Lister:   &EFSBackupPolicyLister{},
	})
}

type EFSBackupPolicyLister struct {
	svc EFSV2Client
}

func (l *EFSBackupPolicyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = efs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := efs.NewDescribeFileSystemsPaginator(svc, &efs.DescribeFileSystemsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range output.FileSystems {
			bpOutput, err := svc.DescribeBackupPolicy(ctx, &efs.DescribeBackupPolicyInput{
				FileSystemId: output.FileSystems[i].FileSystemId,
			})
			if err != nil {
				var notFound *efstypes.PolicyNotFound
				if errors.As(err, &notFound) {
					continue
				}
				return nil, err
			}

			if bpOutput.BackupPolicy == nil {
				continue
			}

			if bpOutput.BackupPolicy.Status != efstypes.StatusEnabled &&
				bpOutput.BackupPolicy.Status != efstypes.StatusEnabling {
				continue
			}

			status := string(bpOutput.BackupPolicy.Status)
			resources = append(resources, &EFSBackupPolicy{
				svc:          svc,
				FileSystemID: output.FileSystems[i].FileSystemId,
				Status:       &status,
			})
		}
	}

	return resources, nil
}

type EFSBackupPolicy struct {
	svc          EFSV2Client
	FileSystemID *string `property:"name=FileSystemId"`
	Status       *string
}

func (r *EFSBackupPolicy) Remove(ctx context.Context) error {
	_, err := r.svc.PutBackupPolicy(ctx, &efs.PutBackupPolicyInput{
		FileSystemId: r.FileSystemID,
		BackupPolicy: &efstypes.BackupPolicy{
			Status: efstypes.StatusDisabled,
		},
	})
	return err
}

func (r *EFSBackupPolicy) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EFSBackupPolicy) String() string {
	return fmt.Sprintf("%s -> BackupPolicy", *r.FileSystemID)
}
