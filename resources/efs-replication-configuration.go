package resources

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EFSReplicationConfigurationResource = "EFSReplicationConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     EFSReplicationConfigurationResource,
		Scope:    nuke.Account,
		Resource: &EFSReplicationConfiguration{},
		Lister:   &EFSReplicationConfigurationLister{},
	})
}

type EFSReplicationConfigurationLister struct {
	svc EFSV2Client
}

func (l *EFSReplicationConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = efs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &efs.DescribeReplicationConfigurationsInput{}
	for {
		resp, err := svc.DescribeReplicationConfigurations(ctx, params)
		if err != nil {
			var notFound *efstypes.ReplicationNotFound
			if errors.As(err, &notFound) {
				return nil, nil
			}
			return nil, err
		}

		for _, repl := range resp.Replications {
			resources = append(resources, &EFSReplicationConfiguration{
				svc:                         svc,
				FileSystemID:                repl.SourceFileSystemId,
				OriginalSourceFileSystemArn: repl.OriginalSourceFileSystemArn,
				SourceFileSystemRegion:      repl.SourceFileSystemRegion,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type EFSReplicationConfiguration struct {
	svc                         EFSV2Client
	FileSystemID                *string
	OriginalSourceFileSystemArn *string
	SourceFileSystemRegion      *string
}

func (r *EFSReplicationConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteReplicationConfiguration(ctx, &efs.DeleteReplicationConfigurationInput{
		SourceFileSystemId: r.FileSystemID,
	})
	return err
}

func (r *EFSReplicationConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EFSReplicationConfiguration) String() string {
	return *r.FileSystemID
}
