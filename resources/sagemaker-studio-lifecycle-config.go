package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerStudioLifecycleConfigResource = "SageMakerStudioLifecycleConfig"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerStudioLifecycleConfigResource,
		Scope:    nuke.Account,
		Resource: &SageMakerStudioLifecycleConfig{},
		Lister:   &SageMakerStudioLifecycleConfigLister{},
	})
}

type SageMakerStudioLifecycleConfigLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerStudioLifecycleConfigLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListStudioLifecycleConfigsPaginator(svc, &sagemaker.ListStudioLifecycleConfigsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.StudioLifecycleConfigs {
			resources = append(resources, &SageMakerStudioLifecycleConfig{
				svc:                       svc,
				StudioLifecycleConfigName: item.StudioLifecycleConfigName,
				StudioLifecycleConfigArn:  item.StudioLifecycleConfigArn,
			})
		}
	}

	return resources, nil
}

type SageMakerStudioLifecycleConfig struct {
	svc                       SageMakerV2Client
	StudioLifecycleConfigName *string
	StudioLifecycleConfigArn  *string
}

func (r *SageMakerStudioLifecycleConfig) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteStudioLifecycleConfig(ctx, &sagemaker.DeleteStudioLifecycleConfigInput{
		StudioLifecycleConfigName: r.StudioLifecycleConfigName,
	})
	return err
}

func (r *SageMakerStudioLifecycleConfig) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerStudioLifecycleConfig) String() string {
	return *r.StudioLifecycleConfigName
}
