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

func Test_Mock_IVSRealtimeIngestConfiguration_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListIngestConfigurations", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListIngestConfigurationsOutput{
			IngestConfigurations: []ivsrealtimetypes.IngestConfigurationSummary{
				{
					Arn:            ptr.String("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123"),
					Name:           ptr.String("my-ingest-config"),
					IngestProtocol: ivsrealtimetypes.IngestProtocolRtmp,
					ParticipantId:  ptr.String("participant-1"),
					StageArn:       ptr.String("arn:aws:ivs:us-east-1:123456789012:stage/stage-1"),
					State:          ivsrealtimetypes.IngestConfigurationStateActive,
				},
			},
		}, nil)

	lister := &IVSRealtimeIngestConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	ingestConfig := resources[0].(*IVSRealtimeIngestConfiguration)
	assertions.Equal("my-ingest-config", *ingestConfig.Name)
	assertions.Equal("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123", *ingestConfig.Arn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeIngestConfiguration_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListIngestConfigurations", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListIngestConfigurationsOutput{
			IngestConfigurations: []ivsrealtimetypes.IngestConfigurationSummary{},
		}, nil)

	lister := &IVSRealtimeIngestConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeIngestConfiguration_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	ingestConfig := &IVSRealtimeIngestConfiguration{
		svc: mockClient,
		Arn: ptr.String("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123"),
	}

	mockClient.On("DeleteIngestConfiguration", mock.Anything, &ivsrealtime.DeleteIngestConfigurationInput{
		Arn: ingestConfig.Arn,
	}).Return(&ivsrealtime.DeleteIngestConfigurationOutput{}, nil)

	assertions.NoError(ingestConfig.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeIngestConfiguration_Properties(t *testing.T) {
	assertions := assert.New(t)

	ingestConfig := IVSRealtimeIngestConfiguration{
		Arn:  ptr.String("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123"),
		Name: ptr.String("my-ingest-config"),
	}

	properties := ingestConfig.Properties()
	assertions.Equal("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123", properties.Get("Arn"))
	assertions.Equal("my-ingest-config", properties.Get("Name"))
}

func Test_Mock_IVSRealtimeIngestConfiguration_String(t *testing.T) {
	assertions := assert.New(t)

	ingestConfig := IVSRealtimeIngestConfiguration{
		Arn: ptr.String("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123"),
	}
	assertions.Equal("arn:aws:ivs:us-east-1:123456789012:ingest-configuration/abc123", ingestConfig.String())
}
