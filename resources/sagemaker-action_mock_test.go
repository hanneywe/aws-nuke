package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

func Test_Mock_SageMakerAction_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListActions", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListActionsOutput{
				ActionSummaries: []sagemakertypes.ActionSummary{
					{
						ActionName: ptr.String("test-action"),
						ActionArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:action/test-action"),
					},
				},
			}, nil,
		)

	lister := &SageMakerActionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	action := resources[0].(*SageMakerAction)
	assertions.Equal("test-action", *action.ActionName)
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:action/test-action", *action.ActionArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerAction_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListActions", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListActionsOutput{
				ActionSummaries: []sagemakertypes.ActionSummary{},
			}, nil,
		)

	lister := &SageMakerActionLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerAction_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	action := &SageMakerAction{
		svc:        mockClient,
		ActionName: ptr.String("test-action"),
		ActionArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:action/test-action"),
	}

	mockClient.
		On("DeleteAction", mock.Anything,
			&sagemaker.DeleteActionInput{
				ActionName: action.ActionName,
			},
		).
		Return(&sagemaker.DeleteActionOutput{}, nil)

	err := action.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerAction_Properties(t *testing.T) {
	assertions := assert.New(t)

	action := SageMakerAction{
		ActionName: ptr.String("test-action"),
		ActionArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:action/test-action"),
	}

	properties := action.Properties()

	assertions.Equal("test-action", properties.Get("ActionName"))
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:action/test-action", properties.Get("ActionArn"))
}

func Test_Mock_SageMakerAction_String(t *testing.T) {
	assertions := assert.New(t)

	action := SageMakerAction{
		ActionName: ptr.String("test-action"),
	}

	assertions.Equal("test-action", action.String())
}
