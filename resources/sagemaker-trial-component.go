package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerTrialComponentResource = "SageMakerTrialComponent"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerTrialComponentResource,
		Scope:    nuke.Account,
		Resource: &SageMakerTrialComponent{},
		Lister:   &SageMakerTrialComponentLister{},
	})
}

type SageMakerTrialComponentLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerTrialComponentLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListTrialComponentsPaginator(svc, &sagemaker.ListTrialComponentsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.TrialComponentSummaries {
			resources = append(resources, &SageMakerTrialComponent{
				svc:                svc,
				TrialComponentName: item.TrialComponentName,
				TrialComponentArn:  item.TrialComponentArn,
			})
		}
	}

	return resources, nil
}

type SageMakerTrialComponent struct {
	svc                SageMakerV2Client
	TrialComponentName *string
	TrialComponentArn  *string
}

func (r *SageMakerTrialComponent) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTrialComponent(ctx, &sagemaker.DeleteTrialComponentInput{
		TrialComponentName: r.TrialComponentName,
	})
	return err
}

func (r *SageMakerTrialComponent) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerTrialComponent) String() string {
	return *r.TrialComponentName
}
