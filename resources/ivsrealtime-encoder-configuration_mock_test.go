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

func Test_Mock_IVSRealtimeEncoderConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListEncoderConfigurations", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListEncoderConfigurationsOutput{
			EncoderConfigurations: []ivsrealtimetypes.EncoderConfigurationSummary{
				{
					Arn:  ptr.String("arn:aws:ivs:us-east-1:123456789012:encoder-configuration/abc123"),
					Name: ptr.String("my-encoder"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &IVSRealtimeEncoderConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ec := resources[0].(*IVSRealtimeEncoderConfiguration)
	a.Equal("my-encoder", *ec.Name)
	a.Equal("test", ec.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeEncoderConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListEncoderConfigurations", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListEncoderConfigurationsOutput{
			EncoderConfigurations: []ivsrealtimetypes.EncoderConfigurationSummary{},
		}, nil)

	lister := &IVSRealtimeEncoderConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeEncoderConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	ec := &IVSRealtimeEncoderConfiguration{
		svc: mockClient,
		ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:encoder-configuration/abc123"),
	}

	mockClient.On("DeleteEncoderConfiguration", mock.Anything, &ivsrealtime.DeleteEncoderConfigurationInput{
		Arn: ec.ARN,
	}).Return(&ivsrealtime.DeleteEncoderConfigurationOutput{}, nil)

	a.NoError(ec.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeEncoderConfiguration_Properties(t *testing.T) {
	a := assert.New(t)

	ec := IVSRealtimeEncoderConfiguration{
		ARN:  ptr.String("arn:aws:ivs:us-east-1:123456789012:encoder-configuration/abc123"),
		Name: ptr.String("my-encoder"),
		Tags: map[string]string{"env": "test"},
	}

	props := ec.Properties()
	a.Equal("my-encoder", props.Get("Name"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_IVSRealtimeEncoderConfiguration_String(t *testing.T) {
	a := assert.New(t)
	ec := IVSRealtimeEncoderConfiguration{ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:encoder-configuration/abc123")}
	a.Equal("arn:aws:ivs:us-east-1:123456789012:encoder-configuration/abc123", ec.String())
}
