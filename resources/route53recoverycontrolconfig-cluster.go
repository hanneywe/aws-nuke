package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53RecoveryControlConfigClusterResource = "Route53RecoveryControlConfigCluster"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53RecoveryControlConfigClusterResource,
		Scope:    nuke.Account,
		Resource: &Route53RecoveryControlConfigCluster{},
		Lister:   &Route53RecoveryControlConfigClusterLister{},
		DependsOn: []string{
			Route53RecoveryControlConfigControlPanelResource,
		},
	})
}

type Route53RecoveryControlConfigClusterLister struct {
	svc Route53RecoveryControlConfigClient
}

func (l *Route53RecoveryControlConfigClusterLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = route53recoverycontrolconfig.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := route53recoverycontrolconfig.NewListClustersPaginator(svc, &route53recoverycontrolconfig.ListClustersInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, c := range resp.Clusters {
			resources = append(resources, &Route53RecoveryControlConfigCluster{
				svc:        svc,
				ClusterArn: c.ClusterArn,
				Name:       c.Name,
			})
		}
	}
	return resources, nil
}

type Route53RecoveryControlConfigCluster struct {
	svc        Route53RecoveryControlConfigClient
	ClusterArn *string
	Name       *string
}

func (r *Route53RecoveryControlConfigCluster) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCluster(ctx, &route53recoverycontrolconfig.DeleteClusterInput{
		ClusterArn: r.ClusterArn,
	})
	return err
}

func (r *Route53RecoveryControlConfigCluster) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53RecoveryControlConfigCluster) String() string {
	return *r.Name
}
