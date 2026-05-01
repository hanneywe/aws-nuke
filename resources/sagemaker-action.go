package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerActionResource = "SageMakerAction"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerActionResource,
		Scope:    nuke.Account,
		Resource: &SageMakerAction{},
		Lister:   &SageMakerActionLister{},
	})
}

type SageMakerActionLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerActionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sagemaker.ListActionsInput{}

	for {
		output, err := svc.ListActions(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, action := range output.ActionSummaries {
			resources = append(resources, &SageMakerAction{
				svc:        svc,
				ActionName: action.ActionName,
				ActionArn:  action.ActionArn,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SageMakerAction struct {
	svc SageMakerV2Client

	ActionName *string
	ActionArn  *string
}

func (r *SageMakerAction) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAction(ctx, &sagemaker.DeleteActionInput{
		ActionName: r.ActionName,
	})
	return err
}

func (r *SageMakerAction) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerAction) String() string {
	return *r.ActionName
}
