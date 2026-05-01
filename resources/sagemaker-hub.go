package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SageMakerHubResource = "SageMakerHub"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerHubResource,
		Scope:    nuke.Account,
		Resource: &SageMakerHub{},
		Lister:   &SageMakerHubLister{},
	})
}

type SageMakerHubLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerHubLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sagemaker.ListHubsInput{}
	for {
		resp, err := svc.ListHubs(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.HubSummaries {
			resources = append(resources, &SageMakerHub{
				svc:     svc,
				HubName: item.HubName,
				HubArn:  item.HubArn,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type SageMakerHub struct {
	svc     SageMakerV2Client
	HubName *string
	HubArn  *string
}

func (r *SageMakerHub) Filter() error {
	if r.HubArn != nil && strings.Contains(*r.HubArn, ":aws:hub/") {
		return fmt.Errorf("cannot delete AWS-managed hub")
	}
	return nil
}

func (r *SageMakerHub) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteHub(ctx, &sagemaker.DeleteHubInput{
		HubName: r.HubName,
	})
	return err
}

func (r *SageMakerHub) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerHub) String() string {
	return *r.HubName
}
