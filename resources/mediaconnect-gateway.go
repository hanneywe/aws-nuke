package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediaconnect"
	mediaconnecttypes "github.com/aws/aws-sdk-go-v2/service/mediaconnect/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MediaConnectGatewayResource = "MediaConnectGateway"

func init() {
	registry.Register(&registry.Registration{
		Name:     MediaConnectGatewayResource,
		Scope:    nuke.Account,
		Resource: &MediaConnectGateway{},
		Lister:   &MediaConnectGatewayLister{},
	})
}

type MediaConnectGatewayLister struct {
	svc MediaConnectClient
}

func (l *MediaConnectGatewayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = mediaconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := mediaconnect.NewListGatewaysPaginator(svc, &mediaconnect.ListGatewaysInput{
		MaxResults: aws.Int32(100),
	})
	for paginator.HasMorePages() {
		listGatewaysOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, gateway := range listGatewaysOutput.Gateways {
			resources = append(resources, &MediaConnectGateway{
				svc:          svc,
				GatewayArn:   gateway.GatewayArn,
				Name:         gateway.Name,
				GatewayState: gateway.GatewayState,
			})
		}
	}

	return resources, nil
}

type MediaConnectGateway struct {
	svc          MediaConnectClient
	GatewayArn   *string
	Name         *string
	GatewayState mediaconnecttypes.GatewayState `property:"-"`
}

func (r *MediaConnectGateway) Filter() error {
	if r.GatewayState == mediaconnecttypes.GatewayStateDeleting ||
		r.GatewayState == mediaconnecttypes.GatewayStateDeleted {
		return fmt.Errorf("already %s", r.GatewayState)
	}
	return nil
}

func (r *MediaConnectGateway) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGateway(ctx, &mediaconnect.DeleteGatewayInput{
		GatewayArn: r.GatewayArn,
	})
	return err
}

func (r *MediaConnectGateway) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MediaConnectGateway) String() string {
	return *r.GatewayArn
}
