package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GameLiftVpcPeeringAuthorizationResource = "GameLiftVpcPeeringAuthorization"

func init() {
	registry.Register(&registry.Registration{
		Name:     GameLiftVpcPeeringAuthorizationResource,
		Scope:    nuke.Account,
		Resource: &GameLiftVpcPeeringAuthorization{},
		Lister:   &GameLiftVpcPeeringAuthorizationLister{},
	})
}

type GameLiftVpcPeeringAuthorizationLister struct {
	svc GameLiftV2Client
}

func (l *GameLiftVpcPeeringAuthorizationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = gamelift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	resp, err := svc.DescribeVpcPeeringAuthorizations(ctx, &gamelift.DescribeVpcPeeringAuthorizationsInput{})
	if err != nil {
		return nil, err
	}
	for _, auth := range resp.VpcPeeringAuthorizations {
		resources = append(resources, &GameLiftVpcPeeringAuthorization{
			svc:                  svc,
			GameLiftAwsAccountID: auth.GameLiftAwsAccountId,
			PeerVpcID:            auth.PeerVpcId,
		})
	}
	return resources, nil
}

type GameLiftVpcPeeringAuthorization struct {
	svc                  GameLiftV2Client
	GameLiftAwsAccountID *string `property:"name=GameLiftAwsAccountId"`
	PeerVpcID            *string `property:"name=PeerVpcId"`
}

func (r *GameLiftVpcPeeringAuthorization) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteVpcPeeringAuthorization(ctx, &gamelift.DeleteVpcPeeringAuthorizationInput{
		GameLiftAwsAccountId: r.GameLiftAwsAccountID,
		PeerVpcId:            r.PeerVpcID,
	})
	return err
}

func (r *GameLiftVpcPeeringAuthorization) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GameLiftVpcPeeringAuthorization) String() string {
	return fmt.Sprintf("%s:%s", *r.GameLiftAwsAccountID, *r.PeerVpcID)
}
