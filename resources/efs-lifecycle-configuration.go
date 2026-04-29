package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EFSLifecycleConfigurationResource = "EFSLifecycleConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     EFSLifecycleConfigurationResource,
		Scope:    nuke.Account,
		Resource: &EFSLifecycleConfiguration{},
		Lister:   &EFSLifecycleConfigurationLister{},
	})
}

type EFSLifecycleConfigurationLister struct {
	svc EFSV2Client
}

func (l *EFSLifecycleConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = efs.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	fsParams := &efs.DescribeFileSystemsInput{}
	for {
		fsResp, err := svc.DescribeFileSystems(ctx, fsParams)
		if err != nil {
			return nil, err
		}
		for i := range fsResp.FileSystems {
			fs := &fsResp.FileSystems[i]
			lcResp, err := svc.DescribeLifecycleConfiguration(ctx, &efs.DescribeLifecycleConfigurationInput{
				FileSystemId: fs.FileSystemId,
			})
			if err != nil {
				return nil, err
			}
			if len(lcResp.LifecyclePolicies) > 0 {
				resources = append(resources, &EFSLifecycleConfiguration{
					svc:          svc,
					FileSystemID: fs.FileSystemId,
				})
			}
		}
		if fsResp.NextMarker == nil {
			break
		}
		fsParams.Marker = fsResp.NextMarker
	}
	return resources, nil
}

type EFSLifecycleConfiguration struct {
	svc          EFSV2Client
	FileSystemID *string
}

func (r *EFSLifecycleConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.PutLifecycleConfiguration(ctx, &efs.PutLifecycleConfigurationInput{
		FileSystemId:      r.FileSystemID,
		LifecyclePolicies: []efstypes.LifecyclePolicy{},
	})
	return err
}

func (r *EFSLifecycleConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EFSLifecycleConfiguration) String() string {
	return *r.FileSystemID
}
