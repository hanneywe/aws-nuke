package resources

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_CloudFrontKeyValueStore_List(t *testing.T) {
	a := assert.New(t)
	mockSvc := new(mockCloudFrontClient)

	now := time.Now()
	mockSvc.On("ListKeyValueStores", mock.Anything, mock.Anything).Return(&cloudfront.ListKeyValueStoresOutput{
		KeyValueStoreList: &types.KeyValueStoreList{
			Items: []types.KeyValueStore{
				{
					Name:             ptr.String("test-kvs"),
					Id:               ptr.String("kvs-123"),
					ARN:              ptr.String("arn:aws:cloudfront::123456789012:key-value-store/test-kvs"),
					Comment:          ptr.String("test comment"),
					Status:           ptr.String("READY"),
					LastModifiedTime: &now,
				},
			},
			Quantity: ptr.Int32(1),
			MaxItems: ptr.Int32(100),
		},
	}, nil)

	lister := &CloudFrontKeyValueStoreLister{svc: mockSvc}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Config: &aws.Config{Region: "us-east-1"},
	})
	a.NoError(err)
	a.Len(resources, 1)

	kvs := resources[0].(*CloudFrontKeyValueStore)
	a.Equal("test-kvs", *kvs.Name)
	a.Equal("kvs-123", *kvs.ID)
}

func Test_Mock_CloudFrontKeyValueStore_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockSvc := new(mockCloudFrontClient)

	mockSvc.On("ListKeyValueStores", mock.Anything, mock.Anything).Return(&cloudfront.ListKeyValueStoresOutput{
		KeyValueStoreList: &types.KeyValueStoreList{
			Items:    []types.KeyValueStore{},
			Quantity: ptr.Int32(0),
			MaxItems: ptr.Int32(100),
		},
	}, nil)

	lister := &CloudFrontKeyValueStoreLister{svc: mockSvc}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{
		Config: &aws.Config{Region: "us-east-1"},
	})
	a.NoError(err)
	a.Len(resources, 0)
}

func Test_Mock_CloudFrontKeyValueStore_Remove(t *testing.T) {
	a := assert.New(t)
	mockSvc := new(mockCloudFrontClient)

	mockSvc.On("DescribeKeyValueStore", mock.Anything, mock.Anything).Return(&cloudfront.DescribeKeyValueStoreOutput{
		ETag: ptr.String("test-etag"),
	}, nil)
	mockSvc.On("DeleteKeyValueStore", mock.Anything, mock.Anything).Return(&cloudfront.DeleteKeyValueStoreOutput{}, nil)

	r := &CloudFrontKeyValueStore{
		svc:  mockSvc,
		Name: ptr.String("test-kvs"),
	}
	err := r.Remove(context.TODO())
	a.NoError(err)
}

func Test_Mock_CloudFrontKeyValueStore_Properties(t *testing.T) {
	a := assert.New(t)

	now := time.Now()
	r := &CloudFrontKeyValueStore{
		Name:             ptr.String("test-kvs"),
		ID:               ptr.String("kvs-123"),
		ARN:              ptr.String("arn:aws:cloudfront::123456789012:key-value-store/test-kvs"),
		Comment:          ptr.String("test comment"),
		Status:           ptr.String("READY"),
		LastModifiedTime: &now,
	}

	props := r.Properties()
	a.Equal("test-kvs", props.Get("Name"))
	a.Equal("kvs-123", props.Get("ID"))
	a.Equal("arn:aws:cloudfront::123456789012:key-value-store/test-kvs", props.Get("ARN"))
	a.Equal("test comment", props.Get("Comment"))
	a.Equal("READY", props.Get("Status"))
}

func Test_Mock_CloudFrontKeyValueStore_String(t *testing.T) {
	a := assert.New(t)

	r := &CloudFrontKeyValueStore{
		Name: ptr.String("test-kvs"),
	}
	a.Equal("test-kvs", r.String())
}
