package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerMlflowAppResource = "SageMakerMlflowApp"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerMlflowAppResource,
		Scope:    nuke.Account,
		Resource: &SageMakerMlflowApp{},
		Lister:   &SageMakerMlflowAppLister{},
	})
}

type SageMakerMlflowAppLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerMlflowAppLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sagemaker.ListMlflowAppsInput{}
	for {
		output, err := svc.ListMlflowApps(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, app := range output.Summaries {
			resources = append(resources, &SageMakerMlflowApp{
				svc:              svc,
				Name:             app.Name,
				ARN:              app.Arn,
				Status:           app.Status,
				MlflowVersion:    app.MlflowVersion,
				CreationTime:     app.CreationTime,
				LastModifiedTime: app.LastModifiedTime,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SageMakerMlflowApp struct {
	svc              SageMakerV2Client
	Name             *string
	ARN              *string
	Status           sagemakertypes.MlflowAppStatus `property:"name=Status"`
	MlflowVersion    *string
	CreationTime     *time.Time
	LastModifiedTime *time.Time
}

func (r *SageMakerMlflowApp) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMlflowApp(ctx, &sagemaker.DeleteMlflowAppInput{
		Arn: r.ARN,
	})
	return err
}

func (r *SageMakerMlflowApp) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerMlflowApp) String() string {
	return *r.Name
}

func (r *SageMakerMlflowApp) Filter() error {
	if r.Status == sagemakertypes.MlflowAppStatusDeleting {
		return fmt.Errorf("already deleting")
	}
	if r.Status == sagemakertypes.MlflowAppStatusDeleted {
		return fmt.Errorf("already deleted")
	}
	return nil
}
