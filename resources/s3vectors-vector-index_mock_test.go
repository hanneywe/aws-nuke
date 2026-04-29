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

func Test_Mock_S3VectorIndex_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3VectorsClient)

	now := time.Now()
	mockClient.On("ListVectorBuckets", mock.Anything, mock.Anything).
		Return(&s3vectors.ListVectorBucketsOutput{
			VectorBuckets: []types.VectorBucketSummary{
				{
					VectorBucketName: ptr.String("my-bucket"),
					VectorBucketArn:  ptr.String("arn:aws:s3vectors:us-east-1:123456789012:bucket/my-bucket"),
					CreationTime:     &now,
				},
			},
		}, nil)

	mockClient.On("ListIndexes", mock.Anything, mock.Anything).
		Return(&s3vectors.ListIndexesOutput{
			Indexes: []types.IndexSummary{
				{
					VectorBucketName: ptr.String("my-bucket"),
					IndexName:        ptr.String("my-index"),
					IndexArn:         ptr.String("arn:aws:s3vectors:us-east-1:123456789012:bucket/my-bucket/index/my-index"),
					CreationTime:     &now,
				},
			},
		}, nil)

	lister := &S3VectorIndexLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testS3VectorsListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*S3VectorIndex)
	a.Equal("my-bucket", *r.VectorBucketName)
	a.Equal("my-index", *r.IndexName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3VectorIndex_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3VectorsClient)

	mockClient.On("ListVectorBuckets", mock.Anything, mock.Anything).
		Return(&s3vectors.ListVectorBucketsOutput{
			VectorBuckets: []types.VectorBucketSummary{},
		}, nil)

	lister := &S3VectorIndexLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testS3VectorsListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3VectorIndex_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockS3VectorsClient)

	r := &S3VectorIndex{
		svc:              mockClient,
		VectorBucketName: ptr.String("my-bucket"),
		IndexName:        ptr.String("my-index"),
	}

	mockClient.On("DeleteIndex", mock.Anything,
		&s3vectors.DeleteIndexInput{
			VectorBucketName: r.VectorBucketName,
			IndexName:        r.IndexName,
		}).Return(&s3vectors.DeleteIndexOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_S3VectorIndex_Properties(t *testing.T) {
	a := assert.New(t)
	now := time.Now()
	r := &S3VectorIndex{
		VectorBucketName: ptr.String("my-bucket"),
		IndexName:        ptr.String("my-index"),
		IndexARN:         ptr.String("arn:aws:s3vectors:us-east-1:123456789012:bucket/my-bucket/index/my-index"),
		CreationTime:     &now,
	}
	props := r.Properties()
	a.Equal("my-bucket", props.Get("VectorBucketName"))
	a.Equal("my-index", props.Get("IndexName"))
}

func Test_Mock_S3VectorIndex_String(t *testing.T) {
	a := assert.New(t)
	r := &S3VectorIndex{
		VectorBucketName: ptr.String("my-bucket"),
		IndexName:        ptr.String("my-index"),
	}
	a.Equal("my-bucket/my-index", r.String())
}
