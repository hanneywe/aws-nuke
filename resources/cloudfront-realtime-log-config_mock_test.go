package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	cloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func Test_Mock_CloudFrontRealtimeLogConfig_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudFrontClient)

	mockClient.On("ListRealtimeLogConfigs", mock.Anything, mock.Anything).
		Return(&cloudfront.ListRealtimeLogConfigsOutput{
			RealtimeLogConfigs: &cloudfronttypes.RealtimeLogConfigs{
				Items: []cloudfronttypes.RealtimeLogConfig{
					{
						Name:         ptr.String("test-value"),
						ARN:          ptr.String("test-arn"),
						SamplingRate: ptr.Int64(100),
						Fields:       []string{"timestamp"},
						EndPoints:    []cloudfronttypes.EndPoint{},
					},
				},
				IsTruncated: ptr.Bool(false),
				Marker:      ptr.String(""),
				MaxItems:    ptr.Int32(100),
			},
		}, nil)

	lister := &CloudFrontRealtimeLogConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudFrontListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CloudFrontRealtimeLogConfig)
	a.Equal("test-value", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudFrontRealtimeLogConfig_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudFrontClient)

	mockClient.On("ListRealtimeLogConfigs", mock.Anything, mock.Anything).
		Return(&cloudfront.ListRealtimeLogConfigsOutput{
			RealtimeLogConfigs: &cloudfronttypes.RealtimeLogConfigs{
				Items:       []cloudfronttypes.RealtimeLogConfig{},
				IsTruncated: ptr.Bool(false),
				Marker:      ptr.String(""),
				MaxItems:    ptr.Int32(100),
			},
		}, nil)

	lister := &CloudFrontRealtimeLogConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudFrontListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudFrontRealtimeLogConfig_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudFrontClient)

	r := &CloudFrontRealtimeLogConfig{
		svc: mockClient,
		ARN: ptr.String("test-arn"),
	}

	mockClient.On("DeleteRealtimeLogConfig", mock.Anything,
		&cloudfront.DeleteRealtimeLogConfigInput{
			ARN: r.ARN,
		}).Return(&cloudfront.DeleteRealtimeLogConfigOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudFrontRealtimeLogConfig_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CloudFrontRealtimeLogConfig{
		Name: ptr.String("test-name"),
		ARN:  ptr.String("test-arn"),
	}
	props := r.Properties()
	a.Equal("test-name", props.Get("Name"))
	a.Equal("test-arn", props.Get("ARN"))
}

func Test_Mock_CloudFrontRealtimeLogConfig_String(t *testing.T) {
	a := assert.New(t)
	r := &CloudFrontRealtimeLogConfig{
		Name: ptr.String("test-name"),
	}
	a.Equal("test-name", r.String())
}
