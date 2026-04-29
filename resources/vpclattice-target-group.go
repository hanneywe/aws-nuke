package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	latticetypes "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeTargetGroupResource = "VPCLatticeTargetGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeTargetGroupResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeTargetGroup{},
		Lister:   &VPCLatticeTargetGroupLister{},
	})
}

type VPCLatticeTargetGroupLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeTargetGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := vpclattice.NewListTargetGroupsPaginator(svc, &vpclattice.ListTargetGroupsInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Items {
			var tags map[string]string
			if resp.Items[i].Arn != nil {
				tagsResp, err := svc.ListTagsForResource(ctx, &vpclattice.ListTagsForResourceInput{
					ResourceArn: resp.Items[i].Arn,
				})
				if err != nil {
					opts.Logger.Warnf("unable to fetch tags for target group: %s", *resp.Items[i].Arn)
				} else {
					tags = tagsResp.Tags
				}
			}

			resources = append(resources, &VPCLatticeTargetGroup{
				svc:  svc,
				ID:   resp.Items[i].Id,
				ARN:  resp.Items[i].Arn,
				Name: resp.Items[i].Name,
				Type: aws.String(string(resp.Items[i].Type)),
				Tags: tags,
			})
		}
	}

	return resources, nil
}

type VPCLatticeTargetGroup struct {
	svc  VPCLatticeClient
	ID   *string
	ARN  *string
	Name *string
	Type *string
	Tags map[string]string
}

func (r *VPCLatticeTargetGroup) Remove(ctx context.Context) error {
	// Deregister all targets before deleting the target group
	targetsResp, err := r.svc.ListTargets(ctx, &vpclattice.ListTargetsInput{
		TargetGroupIdentifier: r.ARN,
	})
	if err == nil && len(targetsResp.Items) > 0 {
		var targets []latticetypes.Target
		for _, t := range targetsResp.Items {
			targets = append(targets, latticetypes.Target{
				Id:   t.Id,
				Port: t.Port,
			})
		}
		_, err = r.svc.DeregisterTargets(ctx, &vpclattice.DeregisterTargetsInput{
			TargetGroupIdentifier: r.ARN,
			Targets:               targets,
		})
		if err != nil {
			return err
		}
	}

	_, err = r.svc.DeleteTargetGroup(ctx, &vpclattice.DeleteTargetGroupInput{
		TargetGroupIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeTargetGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeTargetGroup) String() string {
	return *r.Name
}
