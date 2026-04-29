package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3tables"
	"github.com/aws/aws-sdk-go-v2/service/s3tables/types"
)

func Test_Mock_S3TableBucket_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3TablesClient)

	now := time.Now()
	mockClient.On("ListTableBuckets", mock.Anything, mock.Anything).
		Return(&s3tables.ListTableBucketsOutput{
			TableBuckets: []types.TableBucketSummary{
				{
					Name:      ptr.String("my-table-bucket"),
					Arn:       ptr.String("arn:aws:s3tables:us-east-1:123456789012:bucket/my-table-bucket"),
					CreatedAt: &now,
				},
			},
		}, nil)

	lister := &S3TableBucketLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*S3TableBucket)
	a.Equal("my-table-bucket", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3TableBucket_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3TablesClient)

	mockClient.On("ListTableBuckets", mock.Anything, mock.Anything).
		Return(&s3tables.ListTableBucketsOutput{
			TableBuckets: []types.TableBucketSummary{},
		}, nil)

	lister := &S3TableBucketLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3TableBucket_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3TablesClient)

	r := &S3TableBucket{
		svc:  mockClient,
		Name: ptr.String("my-table-bucket"),
		ARN:  ptr.String("arn:aws:s3tables:us-east-1:123456789012:bucket/my-table-bucket"),
	}

	mockClient.On("DeleteTableBucket", mock.Anything,
		&s3tables.DeleteTableBucketInput{
			TableBucketARN: r.ARN,
		}).Return(&s3tables.DeleteTableBucketOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3TableBucket_Properties(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	r := &S3TableBucket{
		Name:      ptr.String("my-table-bucket"),
		ARN:       ptr.String("arn:aws:s3tables:us-east-1:123456789012:bucket/my-table-bucket"),
		CreatedAt: &now,
	}
	props := r.Properties()
	a.Equal("my-table-bucket", props.Get("Name"))
	a.Equal("arn:aws:s3tables:us-east-1:123456789012:bucket/my-table-bucket", props.Get("ARN"))
}

func Test_Mock_S3TableBucket_String(t *testing.T) {
	a := assert.New(t)
	r := &S3TableBucket{
		Name: ptr.String("my-table-bucket"),
	}
	a.Equal("my-table-bucket", r.String())
}
