package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"
	bedrockagentruntimetypes "github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime/types"
)

func Test_Mock_BedrockAgentRuntimeSession_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBedrockagentruntimeClient)

	now := time.Now()
	mockClient.On("ListSessions", mock.Anything, mock.Anything).
		Return(&bedrockagentruntime.ListSessionsOutput{
			SessionSummaries: []bedrockagentruntimetypes.SessionSummary{
				{
					SessionId:     ptr.String("test-session-id"),
					SessionArn:    ptr.String("arn:aws:bedrock:us-east-1:123456789012:session/test-session-id"),
					SessionStatus: bedrockagentruntimetypes.SessionStatusActive,
					CreatedAt:     &now,
					LastUpdatedAt: &now,
				},
			},
		}, nil)

	lister := &BedrockAgentRuntimeSessionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBedrockagentruntimeListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*BedrockAgentRuntimeSession)
	a.Equal("test-session-id", *r.SessionID)
	a.Equal("ACTIVE", *r.SessionStatus)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BedrockAgentRuntimeSession_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBedrockagentruntimeClient)

	mockClient.On("ListSessions", mock.Anything, mock.Anything).
		Return(&bedrockagentruntime.ListSessionsOutput{
			SessionSummaries: []bedrockagentruntimetypes.SessionSummary{},
		}, nil)

	lister := &BedrockAgentRuntimeSessionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBedrockagentruntimeListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BedrockAgentRuntimeSession_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBedrockagentruntimeClient)

	r := &BedrockAgentRuntimeSession{
		svc:       mockClient,
		SessionID: ptr.String("test-session-id"),
	}

	mockClient.On("DeleteSession", mock.Anything,
		&bedrockagentruntime.DeleteSessionInput{
			SessionIdentifier: r.SessionID,
		}).Return(&bedrockagentruntime.DeleteSessionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_BedrockAgentRuntimeSession_Properties(t *testing.T) {
	a := assert.New(t)
	r := &BedrockAgentRuntimeSession{
		SessionID:     ptr.String("test-session-id"),
		SessionArn:    ptr.String("test-session-arn"),
		SessionStatus: ptr.String("ACTIVE"),
	}
	props := r.Properties()
	a.Equal("test-session-id", props.Get("SessionID"))
	a.Equal("test-session-arn", props.Get("SessionArn"))
	a.Equal("ACTIVE", props.Get("SessionStatus"))
}

func Test_Mock_BedrockAgentRuntimeSession_String(t *testing.T) {
	a := assert.New(t)
	r := &BedrockAgentRuntimeSession{
		SessionID: ptr.String("test-session-id"),
	}
	a.Equal("test-session-id", r.String())
}
