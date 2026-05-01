package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerModelPackageGroupResource = "SageMakerModelPackageGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerModelPackageGroupResource,
		Scope:    nuke.Account,
		Resource: &SageMakerModelPackageGroup{},
		Lister:   &SageMakerModelPackageGroupLister{},
	})
}

type SageMakerModelPackageGroupLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerModelPackageGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListModelPackageGroupsPaginator(svc, &sagemaker.ListModelPackageGroupsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ModelPackageGroupSummaryList {
			resources = append(resources, &SageMakerModelPackageGroup{
				svc:                   svc,
				ModelPackageGroupName: item.ModelPackageGroupName,
				ModelPackageGroupArn:  item.ModelPackageGroupArn,
			})
		}
	}

	return resources, nil
}

type SageMakerModelPackageGroup struct {
	svc                   SageMakerV2Client
	ModelPackageGroupName *string
	ModelPackageGroupArn  *string
}

func (r *SageMakerModelPackageGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteModelPackageGroup(ctx, &sagemaker.DeleteModelPackageGroupInput{
		ModelPackageGroupName: r.ModelPackageGroupName,
	})
	return err
}

func (r *SageMakerModelPackageGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerModelPackageGroup) String() string {
	return *r.ModelPackageGroupName
}
