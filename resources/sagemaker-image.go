package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerImageResource = "SageMakerImage"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerImageResource,
		Scope:    nuke.Account,
		Resource: &SageMakerImage{},
		Lister:   &SageMakerImageLister{},
	})
}

type SageMakerImageLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerImageLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sagemaker.ListImagesInput{}

	for {
		output, err := svc.ListImages(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, image := range output.Images {
			resources = append(resources, &SageMakerImage{
				svc:         svc,
				ImageName:   image.ImageName,
				ImageArn:    image.ImageArn,
				ImageStatus: image.ImageStatus,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SageMakerImage struct {
	svc SageMakerV2Client

	ImageName   *string
	ImageArn    *string
	ImageStatus sagemakertypes.ImageStatus
}

func (r *SageMakerImage) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteImage(ctx, &sagemaker.DeleteImageInput{
		ImageName: r.ImageName,
	})
	return err
}

func (r *SageMakerImage) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerImage) String() string {
	return *r.ImageName
}

func (r *SageMakerImage) Filter() error {
	if r.ImageStatus == sagemakertypes.ImageStatusDeleting {
		return fmt.Errorf("already deleting")
	}
	return nil
}
