package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerModelCardResource = "SageMakerModelCard"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerModelCardResource,
		Scope:    nuke.Account,
		Resource: &SageMakerModelCard{},
		Lister:   &SageMakerModelCardLister{},
	})
}

type SageMakerModelCardLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerModelCardLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListModelCardsPaginator(svc, &sagemaker.ListModelCardsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ModelCardSummaries {
			resources = append(resources, &SageMakerModelCard{
				svc:           svc,
				ModelCardName: item.ModelCardName,
				ModelCardArn:  item.ModelCardArn,
			})
		}
	}

	return resources, nil
}

type SageMakerModelCard struct {
	svc           SageMakerV2Client
	ModelCardName *string
	ModelCardArn  *string
}

func (r *SageMakerModelCard) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteModelCard(ctx, &sagemaker.DeleteModelCardInput{
		ModelCardName: r.ModelCardName,
	})
	return err
}

func (r *SageMakerModelCard) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerModelCard) String() string {
	return *r.ModelCardName
}
