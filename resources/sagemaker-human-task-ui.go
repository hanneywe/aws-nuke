package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerHumanTaskUIResource = "SageMakerHumanTaskUi"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerHumanTaskUIResource,
		Scope:    nuke.Account,
		Resource: &SageMakerHumanTaskUI{},
		Lister:   &SageMakerHumanTaskUILister{},
	})
}

type SageMakerHumanTaskUILister struct {
	svc SageMakerV2Client
}

func (l *SageMakerHumanTaskUILister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListHumanTaskUisPaginator(svc, &sagemaker.ListHumanTaskUisInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.HumanTaskUiSummaries {
			resources = append(resources, &SageMakerHumanTaskUI{
				svc:             svc,
				HumanTaskUIName: item.HumanTaskUiName,
				HumanTaskUIArn:  item.HumanTaskUiArn,
			})
		}
	}

	return resources, nil
}

type SageMakerHumanTaskUI struct {
	svc             SageMakerV2Client
	HumanTaskUIName *string
	HumanTaskUIArn  *string
}

func (r *SageMakerHumanTaskUI) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHumanTaskUi(ctx, &sagemaker.DeleteHumanTaskUiInput{
		HumanTaskUiName: r.HumanTaskUIName,
	})
	return err
}

func (r *SageMakerHumanTaskUI) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerHumanTaskUI) String() string {
	return *r.HumanTaskUIName
}
