package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/s3vectors"
	"github.com/aws/aws-sdk-go-v2/service/s3vectors/types"
)

func Test_Mock_S3VectorBucket_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3VectorsClient)

	now := time.Now()
	mockClient.On("ListVectorBuckets", mock.Anything, mock.Anything).
		Return(&s3vectors.ListVectorBucketsOutput{
			VectorBuckets: []types.VectorBucketSummary{
				{
					VectorBucketName: ptr.String("my-vector-bucket"),
					VectorBucketArn:  ptr.String("arn:aws:s3vectors:us-east-1:123456789012:vector-bucket/my-vector-bucket"),
					CreationTime:     &now,
				},
			},
		}, nil)

	lister := &S3VectorBucketLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testS3VectorsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*S3VectorBucket)
	a.Equal("my-vector-bucket", *r.Name)
	a.Equal("arn:aws:s3vectors:us-east-1:123456789012:vector-bucket/my-vector-bucket", *r.ARN)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3VectorBucket_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3VectorsClient)

	mockClient.On("ListVectorBuckets", mock.Anything, mock.Anything).
		Return(&s3vectors.ListVectorBucketsOutput{
			VectorBuckets: []types.VectorBucketSummary{},
		}, nil)

	lister := &S3VectorBucketLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testS3VectorsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3VectorBucket_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3VectorsClient)

	r := &S3VectorBucket{
		svc:  mockClient,
		Name: ptr.String("my-vector-bucket"),
	}

	mockClient.On("DeleteVectorBucket", mock.Anything,
		&s3vectors.DeleteVectorBucketInput{
			VectorBucketName: r.Name,
		}).Return(&s3vectors.DeleteVectorBucketOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3VectorBucket_Properties(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	r := &S3VectorBucket{
		Name:         ptr.String("my-vector-bucket"),
		ARN:          ptr.String("arn:aws:s3vectors:us-east-1:123456789012:vector-bucket/my-vector-bucket"),
		CreationTime: &now,
	}
	props := r.Properties()
	a.Equal("my-vector-bucket", props.Get("Name"))
	a.Equal("arn:aws:s3vectors:us-east-1:123456789012:vector-bucket/my-vector-bucket", props.Get("ARN"))
}

func Test_Mock_S3VectorBucket_String(t *testing.T) {
	a := assert.New(t)
	r := &S3VectorBucket{
		Name: ptr.String("my-vector-bucket"),
	}
	a.Equal("my-vector-bucket", r.String())
}
