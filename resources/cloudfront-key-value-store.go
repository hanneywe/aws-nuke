package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudFrontKeyValueStoreResource = "CloudFrontKeyValueStore"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudFrontKeyValueStoreResource,
		Scope:    nuke.Account,
		Resource: &CloudFrontKeyValueStore{},
		Lister:   &CloudFrontKeyValueStoreLister{},
		DependsOn: []string{
			CloudFrontFunctionResource,
		},
	})
}

type CloudFrontKeyValueStoreLister struct {
	svc CloudFrontClient
}

func (l *CloudFrontKeyValueStoreLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = cloudfront.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &cloudfront.ListKeyValueStoresInput{}

	for {
		resp, err := svc.ListKeyValueStores(ctx, params)
		if err != nil {
			return nil, err
		}

		if resp.KeyValueStoreList != nil {
			for _, item := range resp.KeyValueStoreList.Items {
				resources = append(resources, &CloudFrontKeyValueStore{
					svc:              svc,
					Name:             item.Name,
					ID:               item.Id,
					ARN:              item.ARN,
					Comment:          item.Comment,
					Status:           item.Status,
					LastModifiedTime: item.LastModifiedTime,
				})
			}

			if resp.KeyValueStoreList.NextMarker == nil {
				break
			}
			params.Marker = resp.KeyValueStoreList.NextMarker
		} else {
			break
		}
	}

	return resources, nil
}

type CloudFrontKeyValueStore struct {
	svc              CloudFrontClient
	Name             *string
	ID               *string
	ARN              *string
	Comment          *string
	Status           *string
	LastModifiedTime *time.Time
}

func (r *CloudFrontKeyValueStore) Remove(ctx context.Context) error {
	resp, err := r.svc.DescribeKeyValueStore(ctx, &cloudfront.DescribeKeyValueStoreInput{
		Name: r.Name,
	})
	if err != nil {
		return err
	}

	_, err = r.svc.DeleteKeyValueStore(ctx, &cloudfront.DeleteKeyValueStoreInput{
		Name:    r.Name,
		IfMatch: resp.ETag,
	})
	return err
}

func (r *CloudFrontKeyValueStore) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudFrontKeyValueStore) String() string {
	return *r.Name
}
