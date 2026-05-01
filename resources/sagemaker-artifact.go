package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerArtifactResource = "SageMakerArtifact"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerArtifactResource,
		Scope:    nuke.Account,
		Resource: &SageMakerArtifact{},
		Lister:   &SageMakerArtifactLister{},
	})
}

type SageMakerArtifactLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerArtifactLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := sagemaker.NewListArtifactsPaginator(svc, &sagemaker.ListArtifactsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.ArtifactSummaries {
			resources = append(resources, &SageMakerArtifact{
				svc:         svc,
				ArtifactArn: item.ArtifactArn,
			})
		}
	}

	return resources, nil
}

type SageMakerArtifact struct {
	svc         SageMakerV2Client
	ArtifactArn *string
}

func (r *SageMakerArtifact) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteArtifact(ctx, &sagemaker.DeleteArtifactInput{
		ArtifactArn: r.ArtifactArn,
	})
	return err
}

func (r *SageMakerArtifact) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerArtifact) String() string {
	return *r.ArtifactArn
}
