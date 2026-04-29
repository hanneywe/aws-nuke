package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LightsailRelationalDatabaseResource = "LightsailRelationalDatabase"

func init() {
	registry.Register(&registry.Registration{
		Name:     LightsailRelationalDatabaseResource,
		Scope:    nuke.Account,
		Resource: &LightsailRelationalDatabase{},
		Lister:   &LightsailRelationalDatabaseLister{},
	})
}

type LightsailRelationalDatabaseLister struct {
	svc LightsailClient
}

func (l *LightsailRelationalDatabaseLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = lightsail.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	resp, err := svc.GetRelationalDatabases(ctx, &lightsail.GetRelationalDatabasesInput{})
	if err != nil {
		return nil, err
	}

	for i := range resp.RelationalDatabases {
		item := &resp.RelationalDatabases[i]
		resources = append(resources, &LightsailRelationalDatabase{
			svc:                    svc,
			RelationalDatabaseName: item.Name,
			Engine:                 item.Engine,
			EngineVersion:          item.EngineVersion,
		})
	}

	return resources, nil
}

type LightsailRelationalDatabase struct {
	svc                    LightsailClient
	RelationalDatabaseName *string
	Engine                 *string
	EngineVersion          *string
}

func (r *LightsailRelationalDatabase) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRelationalDatabase(ctx, &lightsail.DeleteRelationalDatabaseInput{
		RelationalDatabaseName: r.RelationalDatabaseName,
		SkipFinalSnapshot:      aws.Bool(true),
	})
	return err
}

func (r *LightsailRelationalDatabase) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LightsailRelationalDatabase) String() string {
	return *r.RelationalDatabaseName
}
