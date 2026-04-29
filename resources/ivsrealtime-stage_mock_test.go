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

func Test_Mock_IVSRealtimeStage_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListStages", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListStagesOutput{
			Stages: []ivsrealtimetypes.StageSummary{
				{
					Arn:  ptr.String("arn:aws:ivs:us-east-1:123456789012:stage/abc123"),
					Name: ptr.String("my-stage"),
					Tags: map[string]string{"env": "test"},
				},
			},
		}, nil)

	lister := &IVSRealtimeStageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	stage := resources[0].(*IVSRealtimeStage)
	a.Equal("my-stage", *stage.Name)
	a.Equal("test", stage.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeStage_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	mockClient.On("ListStages", mock.Anything, mock.Anything).
		Return(&ivsrealtime.ListStagesOutput{
			Stages: []ivsrealtimetypes.StageSummary{},
		}, nil)

	lister := &IVSRealtimeStageLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIVSRealtimeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeStage_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIVSRealtimeClient)

	stage := &IVSRealtimeStage{
		svc: mockClient,
		ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:stage/abc123"),
	}

	mockClient.On("DeleteStage", mock.Anything, &ivsrealtime.DeleteStageInput{
		Arn: stage.ARN,
	}).Return(&ivsrealtime.DeleteStageOutput{}, nil)

	a.NoError(stage.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IVSRealtimeStage_Properties(t *testing.T) {
	a := assert.New(t)

	stage := IVSRealtimeStage{
		ARN:  ptr.String("arn:aws:ivs:us-east-1:123456789012:stage/abc123"),
		Name: ptr.String("my-stage"),
		Tags: map[string]string{"env": "test"},
	}

	props := stage.Properties()
	a.Equal("my-stage", props.Get("Name"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_IVSRealtimeStage_String(t *testing.T) {
	a := assert.New(t)
	stage := IVSRealtimeStage{ARN: ptr.String("arn:aws:ivs:us-east-1:123456789012:stage/abc123")}
	a.Equal("arn:aws:ivs:us-east-1:123456789012:stage/abc123", stage.String())
}
