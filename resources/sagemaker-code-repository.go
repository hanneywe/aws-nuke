package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerCodeRepositoryResource = "SageMakerCodeRepository"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerCodeRepositoryResource,
		Scope:    nuke.Account,
		Resource: &SageMakerCodeRepository{},
		Lister:   &SageMakerCodeRepositoryLister{},
	})
}

type SageMakerCodeRepositoryLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerCodeRepositoryLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListCodeRepositoriesPaginator(svc, &sagemaker.ListCodeRepositoriesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.CodeRepositorySummaryList {
			resources = append(resources, &SageMakerCodeRepository{
				svc:                svc,
				CodeRepositoryName: item.CodeRepositoryName,
				CodeRepositoryArn:  item.CodeRepositoryArn,
			})
		}
	}

	return resources, nil
}

type SageMakerCodeRepository struct {
	svc                SageMakerV2Client
	CodeRepositoryName *string
	CodeRepositoryArn  *string
}

func (r *SageMakerCodeRepository) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCodeRepository(ctx, &sagemaker.DeleteCodeRepositoryInput{
		CodeRepositoryName: r.CodeRepositoryName,
	})
	return err
}

func (r *SageMakerCodeRepository) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerCodeRepository) String() string {
	return *r.CodeRepositoryName
}
