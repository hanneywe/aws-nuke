package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/configservice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ConfigServiceAggregationAuthorizationResource = "ConfigServiceAggregationAuthorization"

func init() {
	registry.Register(&registry.Registration{
		Name:     ConfigServiceAggregationAuthorizationResource,
		Scope:    nuke.Account,
		Resource: &ConfigServiceAggregationAuthorization{},
		Lister:   &ConfigServiceAggregationAuthorizationLister{},
	})
}

type ConfigServiceAggregationAuthorizationLister struct {
	svc ConfigServiceClient
}

func (l *ConfigServiceAggregationAuthorizationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = configservice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &configservice.DescribeAggregationAuthorizationsInput{}
	for {
		resp, err := svc.DescribeAggregationAuthorizations(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, aa := range resp.AggregationAuthorizations {
			resources = append(resources, &ConfigServiceAggregationAuthorization{
				svc:                 svc,
				AuthorizedAccountID: aa.AuthorizedAccountId,
				AuthorizedAwsRegion: aa.AuthorizedAwsRegion,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ConfigServiceAggregationAuthorization struct {
	svc                 ConfigServiceClient
	AuthorizedAccountID *string `property:"name=AuthorizedAccountId"`
	AuthorizedAwsRegion *string
}

func (r *ConfigServiceAggregationAuthorization) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAggregationAuthorization(ctx, &configservice.DeleteAggregationAuthorizationInput{
		AuthorizedAccountId: r.AuthorizedAccountID,
		AuthorizedAwsRegion: r.AuthorizedAwsRegion,
	})
	return err
}

func (r *ConfigServiceAggregationAuthorization) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ConfigServiceAggregationAuthorization) String() string {
	return fmt.Sprintf("%s:%s", *r.AuthorizedAccountID, *r.AuthorizedAwsRegion)
}
