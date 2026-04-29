package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GameLiftLocationResource = "GameLiftLocation"

func init() {
	registry.Register(&registry.Registration{
		Name:     GameLiftLocationResource,
		Scope:    nuke.Account,
		Resource: &GameLiftLocation{},
		Lister:   &GameLiftLocationLister{},
	})
}

type GameLiftLocationLister struct {
	svc GameLiftV2Client
}

func (l *GameLiftLocationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = gamelift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &gamelift.ListLocationsInput{}
	for {
		resp, err := svc.ListLocations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, loc := range resp.Locations {
			resources = append(resources, &GameLiftLocation{
				svc:          svc,
				LocationName: loc.LocationName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GameLiftLocation struct {
	svc          GameLiftV2Client
	LocationName *string
}

func (r *GameLiftLocation) Filter() error {
	if r.LocationName != nil && !strings.HasPrefix(*r.LocationName, "custom-") {
		return fmt.Errorf("cannot delete AWS-managed location")
	}
	return nil
}

func (r *GameLiftLocation) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLocation(ctx, &gamelift.DeleteLocationInput{
		LocationName: r.LocationName,
	})
	return err
}

func (r *GameLiftLocation) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GameLiftLocation) String() string {
	return *r.LocationName
}
