package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lightsailtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func Test_Mock_LightsailBucketAccessKey_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetBuckets", mock.Anything, mock.Anything).
		Return(&lightsail.GetBucketsOutput{
			Buckets: []lightsailtypes.Bucket{
				{Name: ptr.String("test-bucketname")},
			},
		}, nil)

	mockClient.On("GetBucketAccessKeys", mock.Anything, mock.Anything).
		Return(&lightsail.GetBucketAccessKeysOutput{
			AccessKeys: []lightsailtypes.AccessKey{
				{AccessKeyId: ptr.String("test-accesskeyid")},
			},
		}, nil)

	lister := &LightsailBucketAccessKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*LightsailBucketAccessKey)
	a.Equal("test-bucketname", *r.BucketName)
	a.Equal("test-accesskeyid", *r.AccessKeyID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailBucketAccessKey_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetBuckets", mock.Anything, mock.Anything).
		Return(&lightsail.GetBucketsOutput{
			Buckets: []lightsailtypes.Bucket{},
		}, nil)

	lister := &LightsailBucketAccessKeyLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailBucketAccessKey_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	r := &LightsailBucketAccessKey{
		svc:         mockClient,
		BucketName:  ptr.String("test-bucketname"),
		AccessKeyID: ptr.String("test-accesskeyid"),
	}

	mockClient.On("DeleteBucketAccessKey", mock.Anything,
		&lightsail.DeleteBucketAccessKeyInput{
			BucketName:  r.BucketName,
			AccessKeyId: r.AccessKeyID,
		}).Return(&lightsail.DeleteBucketAccessKeyOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailBucketAccessKey_Properties(t *testing.T) {
	a := assert.New(t)
	r := &LightsailBucketAccessKey{
		BucketName:  ptr.String("test-bucketname"),
		AccessKeyID: ptr.String("test-accesskeyid"),
	}
	props := r.Properties()
	a.Equal("test-bucketname", props.Get("BucketName"))
	a.Equal("test-accesskeyid", props.Get("AccessKeyID"))
}

func Test_Mock_LightsailBucketAccessKey_String(t *testing.T) {
	a := assert.New(t)
	r := &LightsailBucketAccessKey{
		AccessKeyID: ptr.String("test-accesskeyid"),
	}
	a.Equal("test-accesskeyid", r.String())
}
