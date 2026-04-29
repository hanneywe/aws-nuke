package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ECRPullTimeUpdateExclusionResource = "ECRPullTimeUpdateExclusion"

func init() {
	registry.Register(&registry.Registration{
		Name:     ECRPullTimeUpdateExclusionResource,
		Scope:    nuke.Account,
		Resource: &ECRPullTimeUpdateExclusion{},
		Lister:   &ECRPullTimeUpdateExclusionLister{},
	})
}

type ECRPullTimeUpdateExclusionLister struct {
	svc ECRv2Client
}

func (l *ECRPullTimeUpdateExclusionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ecr.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &ecr.ListPullTimeUpdateExclusionsInput{}
	for {
		resp, err := svc.ListPullTimeUpdateExclusions(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, principalArn := range resp.PullTimeUpdateExclusions {
			arn := principalArn
			resources = append(resources, &ECRPullTimeUpdateExclusion{
				svc:          svc,
				PrincipalArn: &arn,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ECRPullTimeUpdateExclusion struct {
	svc          ECRv2Client
	PrincipalArn *string
}

func (r *ECRPullTimeUpdateExclusion) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterPullTimeUpdateExclusion(ctx, &ecr.DeregisterPullTimeUpdateExclusionInput{
		PrincipalArn: r.PrincipalArn,
	})
	return err
}

func (r *ECRPullTimeUpdateExclusion) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ECRPullTimeUpdateExclusion) String() string {
	return *r.PrincipalArn
}
