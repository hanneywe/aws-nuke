package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/efs"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EFSTagsResource = "EFSTags"

func init() {
	registry.Register(&registry.Registration{
		Name:     EFSTagsResource,
		Scope:    nuke.Account,
		Resource: &EFSTags{},
		Lister:   &EFSTagsLister{},
	})
}

type EFSTagsLister struct {
	svc EFSV2Client
}

func (l *EFSTagsLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			tagResp, err := svc.ListTagsForResource(ctx, &efs.ListTagsForResourceInput{
				ResourceId: fs.FileSystemId,
			})
			if err != nil {
				return nil, err
			}
			var userTagKeys []string
			for _, tag := range tagResp.Tags {
				if tag.Key != nil && !strings.HasPrefix(*tag.Key, "aws:") {
					userTagKeys = append(userTagKeys, *tag.Key)
				}
			}
			if len(userTagKeys) > 0 {
				tagCount := len(userTagKeys)
				resources = append(resources, &EFSTags{
					svc:          svc,
					FileSystemID: fs.FileSystemId,
					TagCount:     &tagCount,
					tagKeys:      userTagKeys,
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

type EFSTags struct {
	svc          EFSV2Client
	FileSystemID *string
	TagCount     *int
	tagKeys      []string
}

func (r *EFSTags) Remove(ctx context.Context) error {
	_, err := r.svc.UntagResource(ctx, &efs.UntagResourceInput{
		ResourceId: r.FileSystemID,
		TagKeys:    r.tagKeys,
	})
	return err
}

func (r *EFSTags) Properties() types.Properties {
	props := types.NewProperties()
	props.Set("FileSystemID", r.FileSystemID)
	if r.TagCount != nil {
		props.Set("TagCount", fmt.Sprintf("%d", *r.TagCount))
	}
	return props
}

func (r *EFSTags) String() string {
	return *r.FileSystemID
}
