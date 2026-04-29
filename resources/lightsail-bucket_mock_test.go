package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lstypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func Test_Mock_LightsailBucket_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetBuckets", mock.Anything, mock.Anything).
		Return(&lightsail.GetBucketsOutput{
			Buckets: []lstypes.Bucket{
				{
					Name: ptr.String("my-bucket"),
					Arn:  ptr.String("arn:aws:lightsail:us-east-1:123456789012:Bucket/my-bucket"),
				},
			},
		}, nil)

	lister := &LightsailBucketLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	b := resources[0].(*LightsailBucket)
	a.Equal("my-bucket", *b.BucketName)
	a.Equal("arn:aws:lightsail:us-east-1:123456789012:Bucket/my-bucket", *b.Arn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailBucket_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetBuckets", mock.Anything, mock.Anything).
		Return(&lightsail.GetBucketsOutput{
			Buckets: []lstypes.Bucket{},
		}, nil)

	lister := &LightsailBucketLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailBucket_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	b := &LightsailBucket{
		svc:        mockClient,
		BucketName: ptr.String("my-bucket"),
	}

	mockClient.On("DeleteBucket", mock.Anything, &lightsail.DeleteBucketInput{
		BucketName:  b.BucketName,
		ForceDelete: ptr.Bool(true),
	}).Return(&lightsail.DeleteBucketOutput{}, nil)

	a.NoError(b.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailBucket_Properties(t *testing.T) {
	a := assert.New(t)

	b := LightsailBucket{
		BucketName: ptr.String("my-bucket"),
		Arn:        ptr.String("arn:aws:lightsail:us-east-1:123456789012:Bucket/my-bucket"),
	}

	props := b.Properties()
	a.Equal("my-bucket", props.Get("BucketName"))
	a.Equal("arn:aws:lightsail:us-east-1:123456789012:Bucket/my-bucket", props.Get("Arn"))
}

func Test_Mock_LightsailBucket_String(t *testing.T) {
	a := assert.New(t)
	b := LightsailBucket{BucketName: ptr.String("my-bucket")}
	a.Equal("my-bucket", b.String())
}
