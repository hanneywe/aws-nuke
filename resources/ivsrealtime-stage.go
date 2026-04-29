package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IVSRealtimeStageResource = "IVSRealtimeStage"

func init() {
	registry.Register(&registry.Registration{
		Name:     IVSRealtimeStageResource,
		Scope:    nuke.Account,
		Resource: &IVSRealtimeStage{},
		Lister:   &IVSRealtimeStageLister{},
	})
}

type IVSRealtimeStageLister struct {
	svc IVSRealtimeClient
}

func (l *IVSRealtimeStageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ivsrealtime.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ivsrealtime.NewListStagesPaginator(svc, &ivsrealtime.ListStagesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, stage := range resp.Stages {
			resources = append(resources, &IVSRealtimeStage{
				svc:  svc,
				ARN:  stage.Arn,
				Name: stage.Name,
				Tags: stage.Tags,
			})
		}
	}

	return resources, nil
}

type IVSRealtimeStage struct {
	svc  IVSRealtimeClient
	ARN  *string
	Name *string
	Tags map[string]string
}

func (r *IVSRealtimeStage) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteStage(ctx, &ivsrealtime.DeleteStageInput{
		Arn: r.ARN,
	})
	return err
}

func (r *IVSRealtimeStage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IVSRealtimeStage) String() string {
	return *r.ARN
}
