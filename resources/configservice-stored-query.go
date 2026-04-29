package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConfigServiceStoredQueryResource = "ConfigServiceStoredQuery"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConfigServiceStoredQueryResource,
		Scope:    nuke.Account,
		Resource: &ConfigServiceStoredQuery{},
		Lister:   &ConfigServiceStoredQueryLister{},
	})
}

type ConfigServiceStoredQueryLister struct {
	svc ConfigServiceClient
}

func (l *ConfigServiceStoredQueryLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = configservice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &configservice.ListStoredQueriesInput{}
	for {
		resp, err := svc.ListStoredQueries(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, q := range resp.StoredQueryMetadata {
			resources = append(resources, &ConfigServiceStoredQuery{
				svc:       svc,
				QueryID:   q.QueryId,
				QueryName: q.QueryName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ConfigServiceStoredQuery struct {
	svc       ConfigServiceClient
	QueryID   *string `property:"name=QueryId"`
	QueryName *string
}

func (r *ConfigServiceStoredQuery) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteStoredQuery(ctx, &configservice.DeleteStoredQueryInput{
		QueryName: r.QueryName,
	})
	return err
}

func (r *ConfigServiceStoredQuery) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConfigServiceStoredQuery) String() string {
	return *r.QueryName
}
