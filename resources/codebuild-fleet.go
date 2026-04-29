package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"
	cbtypes "github.com/aws/aws-sdk-go-v2/service/codebuild/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CodeBuildFleetResource = "CodeBuildFleet"

func init() {
	registry.Register(&registry.Registration{
		Name:     CodeBuildFleetResource,
		Scope:    nuke.Account,
		Resource: &CodeBuildFleet{},
		Lister:   &CodeBuildFleetLister{},
	})
}

type CodeBuildFleetLister struct {
	svc CodeBuildClient
}

func (l *CodeBuildFleetLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = codebuild.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := codebuild.NewListFleetsPaginator(svc, &codebuild.ListFleetsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(resp.Fleets) == 0 {
			continue
		}

		batchResp, err := svc.BatchGetFleets(ctx, &codebuild.BatchGetFleetsInput{
			Names: resp.Fleets,
		})
		if err != nil {
			return nil, err
		}

		for i := range batchResp.Fleets {
			fleet := &batchResp.Fleets[i]
			tags := make(map[string]string)
			for _, t := range fleet.Tags {
				if t.Key != nil && t.Value != nil {
					tags[*t.Key] = *t.Value
				}
			}

			var status cbtypes.FleetStatusCode
			if fleet.Status != nil {
				status = fleet.Status.StatusCode
			}

			resources = append(resources, &CodeBuildFleet{
				svc:    svc,
				Name:   fleet.Name,
				ARN:    fleet.Arn,
				Status: status,
				Tags:   tags,
			})
		}
	}

	return resources, nil
}

type CodeBuildFleet struct {
	svc    CodeBuildClient
	Name   *string
	ARN    *string
	Status cbtypes.FleetStatusCode
	Tags   map[string]string
}

func (r *CodeBuildFleet) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFleet(ctx, &codebuild.DeleteFleetInput{
		Arn: r.ARN,
	})
	return err
}

func (r *CodeBuildFleet) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CodeBuildFleet) String() string {
	return *r.Name
}
