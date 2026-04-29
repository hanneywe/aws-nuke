package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2NetworkInsightsAccessScopeResource = "EC2NetworkInsightsAccessScope"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2NetworkInsightsAccessScopeResource,
		Scope:    nuke.Account,
		Resource: &EC2NetworkInsightsAccessScope{},
		Lister:   &EC2NetworkInsightsAccessScopeLister{},
	})
}

type EC2NetworkInsightsAccessScopeLister struct {
	svc EC2Client
}

func (l *EC2NetworkInsightsAccessScopeLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeNetworkInsightsAccessScopesPaginator(svc,
		&ec2.DescribeNetworkInsightsAccessScopesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, scope := range resp.NetworkInsightsAccessScopes {
			resources = append(resources, &EC2NetworkInsightsAccessScope{
				svc:                           svc,
				NetworkInsightsAccessScopeID:  scope.NetworkInsightsAccessScopeId,
				NetworkInsightsAccessScopeArn: scope.NetworkInsightsAccessScopeArn,
				Tags:                          scope.Tags,
			})
		}
	}

	return resources, nil
}

type EC2NetworkInsightsAccessScope struct {
	svc                           EC2Client
	NetworkInsightsAccessScopeID  *string `property:"name=NetworkInsightsAccessScopeId"`
	NetworkInsightsAccessScopeArn *string
	Tags                          []ec2types.Tag
}

func (r *EC2NetworkInsightsAccessScope) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteNetworkInsightsAccessScope(ctx, &ec2.DeleteNetworkInsightsAccessScopeInput{
		NetworkInsightsAccessScopeId: r.NetworkInsightsAccessScopeID,
	})
	return err
}

func (r *EC2NetworkInsightsAccessScope) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2NetworkInsightsAccessScope) String() string {
	return *r.NetworkInsightsAccessScopeID
}
