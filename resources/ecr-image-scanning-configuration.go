package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ECRImageScanningConfigurationResource = "ECRImageScanningConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     ECRImageScanningConfigurationResource,
		Scope:    nuke.Account,
		Resource: &ECRImageScanningConfiguration{},
		Lister:   &ECRImageScanningConfigurationLister{},
	})
}

type ECRImageScanningConfigurationLister struct {
	svc ECRv2Client
}

func (l *ECRImageScanningConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ecr.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &ecr.DescribeRepositoriesInput{}
	for {
		resp, err := svc.DescribeRepositories(ctx, params)
		if err != nil {
			return nil, err
		}
		for i := range resp.Repositories {
			repo := resp.Repositories[i]
			if repo.ImageScanningConfiguration != nil && repo.ImageScanningConfiguration.ScanOnPush {
				resources = append(resources, &ECRImageScanningConfiguration{
					svc:            svc,
					RepositoryName: repo.RepositoryName,
					ScanOnPush:     &repo.ImageScanningConfiguration.ScanOnPush,
				})
			}
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ECRImageScanningConfiguration struct {
	svc            ECRv2Client
	RepositoryName *string
	ScanOnPush     *bool
}

func (r *ECRImageScanningConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.PutImageScanningConfiguration(ctx, &ecr.PutImageScanningConfigurationInput{
		RepositoryName: r.RepositoryName,
		ImageScanningConfiguration: &ecrtypes.ImageScanningConfiguration{
			ScanOnPush: false,
		},
	})
	return err
}

func (r *ECRImageScanningConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ECRImageScanningConfiguration) String() string {
	return *r.RepositoryName
}
