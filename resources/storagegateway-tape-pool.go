package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/storagegateway"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const StorageGatewayTapePoolResource = "StorageGatewayTapePool"

func init() {
	registry.Register(&registry.Registration{
		Name:     StorageGatewayTapePoolResource,
		Scope:    nuke.Account,
		Resource: &StorageGatewayTapePool{},
		Lister:   &StorageGatewayTapePoolLister{},
	})
}

type StorageGatewayTapePoolLister struct {
	svc StorageGatewayV2Client
}

func (l *StorageGatewayTapePoolLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = storagegateway.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &storagegateway.ListTapePoolsInput{}
	for {
		output, err := svc.ListTapePools(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, tapePool := range output.PoolInfos {
			resources = append(resources, &StorageGatewayTapePool{
				svc:      svc,
				PoolARN:  tapePool.PoolARN,
				PoolName: tapePool.PoolName,
			})
		}

		if output.Marker == nil {
			break
		}
		params.Marker = output.Marker
	}

	return resources, nil
}

type StorageGatewayTapePool struct {
	svc      StorageGatewayV2Client
	PoolARN  *string
	PoolName *string
}

func (r *StorageGatewayTapePool) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTapePool(ctx, &storagegateway.DeleteTapePoolInput{
		PoolARN: r.PoolARN,
	})
	return err
}

func (r *StorageGatewayTapePool) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *StorageGatewayTapePool) String() string {
	return *r.PoolARN
}
