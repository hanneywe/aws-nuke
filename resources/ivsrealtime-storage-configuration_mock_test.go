package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ivsrealtime"
	ivsrealtimetypes "github.com/aws/aws-sdk-go-v2/service/ivsrealtime/types"
)

func Test_Mock_IVSRealtimeStorageConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListStorageConfigurations", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListStorageConfigurationsOutput{
			StorageConfigurations: []ivsrealtimetypes.StorageConfigurationSummary{
				{
					Arn:  ptr.String("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123"),
					Name: ptr.String("my-storage-config"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &IVSRealtimeStorageConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*IVSRealtimeStorageConfiguration)
	a.Equal("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123", *r.ARN)
	a.Equal("my-storage-config", *r.Name)
	a.Equal("test", r.Tags["env"])
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeStorageConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListStorageConfigurations", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListStorageConfigurationsOutput{
			StorageConfigurations: []ivsrealtimetypes.StorageConfigurationSummary{},
		}, nil)

	lister := &IVSRealtimeStorageConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeStorageConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	r := &IVSRealtimeStorageConfiguration{
		svc: mockClient,
		ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123"),
	}

	mockClient.On("DeleteStorageConfiguration", mock.Anything,
		&ivsrealtime.DeleteStorageConfigurationInput{
			Arn: r.ARN,
		}).Return(&ivsrealtime.DeleteStorageConfigurationOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeStorageConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := &IVSRealtimeStorageConfiguration{
		ARN:  ptr.String("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123"),
		Name: ptr.String("my-storage-config"),
		Tags: map[string]string{"env": "test"},
	}
	props := r.Properties()
	a.Equal("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123", props.Get("ARN"))
	a.Equal("my-storage-config", props.Get("Name"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_IVSRealtimeStorageConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &IVSRealtimeStorageConfiguration{
		ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123"),
	}
	a.Equal("arn:aws:ivs:us-east-1:123456789012:storage-configuration/abc123", r.String())
}
