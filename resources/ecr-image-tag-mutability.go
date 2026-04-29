package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ECRImageTagMutabilityResource = "ECRImageTagMutability"

func init() {
	registry.Register(&registry.Registration{
		Name:     ECRImageTagMutabilityResource,
		Scope:    nuke.Account,
		Resource: &ECRImageTagMutability{},
		Lister:   &ECRImageTagMutabilityLister{},
	})
}

type ECRImageTagMutabilityLister struct {
	svc ECRv2Client
}

func (l *ECRImageTagMutabilityLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ecr.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &ecr.DescribeRepositoriesInput{}
	for {
		resp, err := svc.DescribeRepositories(ctx, params)
		if err != nil {
			return nil, err
		}
		for i := range resp.Repositories {
			repo := resp.Repositories[i]
			if repo.ImageTagMutability == ecrtypes.ImageTagMutabilityImmutable {
				mutability := string(repo.ImageTagMutability)
				resources = append(resources, &ECRImageTagMutability{
					svc:                svc,
					RepositoryName:     repo.RepositoryName,
					ImageTagMutability: &mutability,
				})
			}
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ECRImageTagMutability struct {
	svc                ECRv2Client
	RepositoryName     *string
	ImageTagMutability *string
}

func (r *ECRImageTagMutability) Remove(ctx context.Context) error {
	_, err := r.svc.PutImageTagMutability(ctx, &ecr.PutImageTagMutabilityInput{
		RepositoryName:     r.RepositoryName,
		ImageTagMutability: ecrtypes.ImageTagMutabilityMutable,
	})
	return err
}

func (r *ECRImageTagMutability) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ECRImageTagMutability) String() string {
	return *r.RepositoryName
}
