package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libTypes "github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ResourceExplorer2SetupResource = "ResourceExplorer2Setup"

func init() {
	registry.Register(&registry.Registration{
		Name:     ResourceExplorer2SetupResource,
		Scope:    nuke.Account,
		Resource: &ResourceExplorer2Setup{},
		Lister:   &ResourceExplorer2SetupLister{},
	})
}

type ResourceExplorer2SetupLister struct {
	svc ResourceExplorer2Client
}

func (l *ResourceExplorer2SetupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			resources = append(resources, &ResourceExplorer2Setup{
				svc:       svc,
				IndexArn:  resp.Indexes[i].Arn,
				IndexType: resp.Indexes[i].Type,
				Region:    resp.Indexes[i].Region,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type ResourceExplorer2Setup struct {
	svc       ResourceExplorer2Client
	IndexArn  *string
	IndexType types.IndexType `property:"name=IndexType"`
	Region    *string
}

func (r *ResourceExplorer2Setup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteIndex(ctx, &resourceexplorer2.DeleteIndexInput{
		Arn: r.IndexArn,
	})
	return err
}

func (r *ResourceExplorer2Setup) Properties() libTypes.Properties {
	return libTypes.NewPropertiesFromStruct(r)
}

func (r *ResourceExplorer2Setup) String() string {
	return *r.IndexArn
}
