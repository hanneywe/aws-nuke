package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/neptune"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NeptuneDBClusterEndpointResource = "NeptuneDBClusterEndpoint"

func init() {
	registry.Register(&registry.Registration{
		Name:     NeptuneDBClusterEndpointResource,
		Scope:    nuke.Account,
		Resource: &NeptuneDBClusterEndpoint{},
		Lister:   &NeptuneDBClusterEndpointLister{},
	})
}

type NeptuneDBClusterEndpointLister struct {
	svc NeptuneV2Client
}

func (l *NeptuneDBClusterEndpointLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = neptune.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &neptune.DescribeDBClusterEndpointsInput{}
	paginator := neptune.NewDescribeDBClusterEndpointsPaginator(svc, params)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, endpoint := range output.DBClusterEndpoints {
			resources = append(resources, &NeptuneDBClusterEndpoint{
				svc:                         svc,
				DBClusterEndpointIdentifier: endpoint.DBClusterEndpointIdentifier,
				DBClusterIdentifier:         endpoint.DBClusterIdentifier,
				EndpointType:                endpoint.EndpointType,
				Status:                      endpoint.Status,
			})
		}
	}

	return resources, nil
}

type NeptuneDBClusterEndpoint struct {
	svc                         NeptuneV2Client
	DBClusterEndpointIdentifier *string
	DBClusterIdentifier         *string
	EndpointType                *string
	Status                      *string
}

func (r *NeptuneDBClusterEndpoint) Filter() error {
	if r.Status != nil && *r.Status == "deleting" {
		return fmt.Errorf("already deleting")
	}
	if r.EndpointType != nil && (*r.EndpointType == "READER" || *r.EndpointType == "WRITER") {
		return fmt.Errorf("default cluster endpoint cannot be deleted")
	}
	if r.DBClusterEndpointIdentifier == nil {
		return fmt.Errorf("no endpoint identifier")
	}
	return nil
}

func (r *NeptuneDBClusterEndpoint) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDBClusterEndpoint(ctx, &neptune.DeleteDBClusterEndpointInput{
		DBClusterEndpointIdentifier: r.DBClusterEndpointIdentifier,
	})
	return err
}

func (r *NeptuneDBClusterEndpoint) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *NeptuneDBClusterEndpoint) String() string {
	if r.DBClusterEndpointIdentifier == nil {
		return ""
	}
	return *r.DBClusterEndpointIdentifier
}
