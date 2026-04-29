package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/gamelift"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const GameLiftAliasResource = "GameLiftAlias"

func init() {
	registry.Register(&registry.Registration{
		Name:     GameLiftAliasResource,
		Scope:    nuke.Account,
		Resource: &GameLiftAlias{},
		Lister:   &GameLiftAliasLister{},
	})
}

type GameLiftAliasLister struct {
	svc GameLiftV2Client
}

func (l *GameLiftAliasLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = gamelift.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &gamelift.ListAliasesInput{}
	for {
		resp, err := svc.ListAliases(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, a := range resp.Aliases {
			resources = append(resources, &GameLiftAlias{
				svc:     svc,
				AliasID: a.AliasId,
				Name:    a.Name,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type GameLiftAlias struct {
	svc     GameLiftV2Client
	AliasID *string `property:"name=AliasId"`
	Name    *string
}

func (r *GameLiftAlias) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAlias(ctx, &gamelift.DeleteAliasInput{
		AliasId: r.AliasID,
	})
	return err
}

func (r *GameLiftAlias) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *GameLiftAlias) String() string {
	return *r.AliasID
}
