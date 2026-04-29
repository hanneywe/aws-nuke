package resources

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libTypes "github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ResourceExplorer2IndexResource = "ResourceExplorer2Index"

func init() {
	registry.Register(&registry.Registration{
		Name:     ResourceExplorer2IndexResource,
		Scope:    nuke.Account,
		Resource: &ResourceExplorer2Index{},
		Lister:   &ResourceExplorer2IndexLister{},
	})
}

type ResourceExplorer2IndexLister struct {
	svc ResourceExplorer2Client
}

func (l *ResourceExplorer2IndexLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = resourceexplorer2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &resourceexplorer2.ListIndexesInput{}
	for {
		resp, err := svc.ListIndexes(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Indexes {
			tags := make(map[string]string)
			tagResp, err := svc.ListTagsForResource(ctx, &resourceexplorer2.ListTagsForResourceInput{
				ResourceArn: resp.Indexes[i].Arn,
			})
			if err != nil {
				logrus.WithError(err).Error("unable to list tags for resource")
			}
			if tagResp != nil {
				tags = tagResp.Tags
			}

			resources = append(resources, &ResourceExplorer2Index{
				svc:    svc,
				ARN:    resp.Indexes[i].Arn,
				Region: resp.Indexes[i].Region,
				Type:   resp.Indexes[i].Type,
				Tags:   tags,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ResourceExplorer2Index struct {
	svc    ResourceExplorer2Client
	ARN    *string
	Region *string
	Type   types.IndexType `property:"name=Type"`
	Tags   map[string]string
}

func (r *ResourceExplorer2Index) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIndex(ctx, &resourceexplorer2.DeleteIndexInput{
		Arn: r.ARN,
	})
	return err
}

func (r *ResourceExplorer2Index) Properties() libTypes.Properties {
	return libTypes.NewPropertiesFromStruct(r)
}

func (r *ResourceExplorer2Index) String() string {
	return *r.ARN
}
