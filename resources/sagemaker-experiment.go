package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerExperimentResource = "SageMakerExperiment"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerExperimentResource,
		Scope:    nuke.Account,
		Resource: &SageMakerExperiment{},
		Lister:   &SageMakerExperimentLister{},
	})
}

type SageMakerExperimentLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerExperimentLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListExperimentsPaginator(svc, &sagemaker.ListExperimentsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ExperimentSummaries {
			resources = append(resources, &SageMakerExperiment{
				svc:            svc,
				ExperimentName: item.ExperimentName,
				ExperimentArn:  item.ExperimentArn,
			})
		}
	}

	return resources, nil
}

type SageMakerExperiment struct {
	svc            SageMakerV2Client
	ExperimentName *string
	ExperimentArn  *string
}

func (r *SageMakerExperiment) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteExperiment(ctx, &sagemaker.DeleteExperimentInput{
		ExperimentName: r.ExperimentName,
	})
	return err
}

func (r *SageMakerExperiment) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerExperiment) String() string {
	return *r.ExperimentName
}
