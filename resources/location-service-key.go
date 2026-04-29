package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/location"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LocationServiceKeyResource = "LocationServiceKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     LocationServiceKeyResource,
		Scope:    nuke.Account,
		Resource: &LocationServiceKey{},
		Lister:   &LocationServiceKeyLister{},
	})
}

type LocationServiceKeyLister struct {
	svc LocationServiceClient
}

func (l *LocationServiceKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = location.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &location.ListKeysInput{}
	for {
		resp, err := svc.ListKeys(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, entry := range resp.Entries {
			resources = append(resources, &LocationServiceKey{
				svc:        svc,
				KeyName:    entry.KeyName,
				CreateTime: entry.CreateTime,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type LocationServiceKey struct {
	svc        LocationServiceClient
	KeyName    *string
	CreateTime *time.Time
}

func (r *LocationServiceKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteKey(ctx, &location.DeleteKeyInput{
		KeyName: r.KeyName,
	})
	return err
}

func (r *LocationServiceKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LocationServiceKey) String() string {
	return *r.KeyName
}
