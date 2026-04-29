package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sfn"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SFNActivityResource = "SFNActivity"

func init() {
	registry.Register(&registry.Registration{
		Name:     SFNActivityResource,
		Scope:    nuke.Account,
		Resource: &SFNActivity{},
		Lister:   &SFNActivityLister{},
	})
}

type SFNActivityLister struct {
	svc SFNv2Client
}

func (l *SFNActivityLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sfn.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sfn.ListActivitiesInput{}
	for {
		resp, err := svc.ListActivities(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Activities {
			resources = append(resources, &SFNActivity{
				svc:         svc,
				Name:        resp.Activities[i].Name,
				ActivityArn: resp.Activities[i].ActivityArn,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type SFNActivity struct {
	svc         SFNv2Client
	Name        *string
	ActivityArn *string
}

func (r *SFNActivity) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteActivity(ctx, &sfn.DeleteActivityInput{
		ActivityArn: r.ActivityArn,
	})
	return err
}

func (r *SFNActivity) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SFNActivity) String() string {
	return *r.Name
}
