package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SESv2DedicatedIPPoolResource = "SESv2DedicatedIpPool"

func init() {
	registry.Register(&registry.Registration{
		Name:     SESv2DedicatedIPPoolResource,
		Scope:    nuke.Account,
		Resource: &SESv2DedicatedIPPool{},
		Lister:   &SESv2DedicatedIPPoolLister{},
	})
}

type SESv2DedicatedIPPoolLister struct {
	svc SESv2Client
}

func (l *SESv2DedicatedIPPoolLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = sesv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := sesv2.NewListDedicatedIpPoolsPaginator(svc, &sesv2.ListDedicatedIpPoolsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, name := range resp.DedicatedIpPools {
			resources = append(resources, &SESv2DedicatedIPPool{
				svc:      svc,
				PoolName: aws.String(name),
			})
		}
	}
	return resources, nil
}

type SESv2DedicatedIPPool struct {
	svc      SESv2Client
	PoolName *string
}

func (r *SESv2DedicatedIPPool) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDedicatedIpPool(ctx, &sesv2.DeleteDedicatedIpPoolInput{
		PoolName: r.PoolName,
	})
	return err
}

func (r *SESv2DedicatedIPPool) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SESv2DedicatedIPPool) String() string {
	return *r.PoolName
}
