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

const SageMakerMlflowTrackingServerResource = "SageMakerMlflowTrackingServer"

func init() {
	registry.Register(&registry.Registration{
		Name:     SageMakerMlflowTrackingServerResource,
		Scope:    nuke.Account,
		Resource: &SageMakerMlflowTrackingServer{},
		Lister:   &SageMakerMlflowTrackingServerLister{},
	})
}

type SageMakerMlflowTrackingServerLister struct {
	svc SageMakerV2Client
}

func (l *SageMakerMlflowTrackingServerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sagemaker.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sagemaker.ListMlflowTrackingServersInput{}

	for {
		output, err := svc.ListMlflowTrackingServers(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, trackingServer := range output.TrackingServerSummaries {
			resources = append(resources, &SageMakerMlflowTrackingServer{
				svc:                  svc,
				TrackingServerName:   trackingServer.TrackingServerName,
				TrackingServerArn:    trackingServer.TrackingServerArn,
				TrackingServerStatus: trackingServer.TrackingServerStatus,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SageMakerMlflowTrackingServer struct {
	svc SageMakerV2Client

	TrackingServerName   *string
	TrackingServerArn    *string
	TrackingServerStatus sagemakertypes.TrackingServerStatus
}

func (r *SageMakerMlflowTrackingServer) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteMlflowTrackingServer(ctx, &sagemaker.DeleteMlflowTrackingServerInput{
		TrackingServerName: r.TrackingServerName,
	})
	return err
}

func (r *SageMakerMlflowTrackingServer) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SageMakerMlflowTrackingServer) String() string {
	return *r.TrackingServerName
}

func (r *SageMakerMlflowTrackingServer) Filter() error {
	if r.TrackingServerStatus == sagemakertypes.TrackingServerStatusDeleting {
		return fmt.Errorf("already deleting")
	}
	return nil
}
