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

func Test_Mock_SageMakerExperiment_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListExperiments", mock.Anything, mock.Anything).
		Return(&sagemaker.ListExperimentsOutput{
			ExperimentSummaries: []sagemakertypes.ExperimentSummary{
				{
					ExperimentName: ptr.String("my-experiment"),
					ExperimentArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:experiment/my-experiment"),
				},
			},
		}, nil)

	lister := &SageMakerExperimentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	experiment := resources[0].(*SageMakerExperiment)
	a.Equal("my-experiment", *experiment.ExperimentName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:experiment/my-experiment", *experiment.ExperimentArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerExperiment_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListExperiments", mock.Anything, mock.Anything).
		Return(&sagemaker.ListExperimentsOutput{
			ExperimentSummaries: []sagemakertypes.ExperimentSummary{},
		}, nil)

	lister := &SageMakerExperimentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerExperiment_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	experiment := &SageMakerExperiment{
		svc:            mockClient,
		ExperimentName: ptr.String("my-experiment"),
	}

	mockClient.On("DeleteExperiment", mock.Anything, &sagemaker.DeleteExperimentInput{
		ExperimentName: experiment.ExperimentName,
	}).Return(&sagemaker.DeleteExperimentOutput{}, nil)

	a.NoError(experiment.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerExperiment_Properties(t *testing.T) {
	a := assert.New(t)

	experiment := SageMakerExperiment{
		ExperimentName: ptr.String("my-experiment"),
		ExperimentArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:experiment/my-experiment"),
	}

	props := experiment.Properties()
	a.Equal("my-experiment", props.Get("ExperimentName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:experiment/my-experiment", props.Get("ExperimentArn"))
}

func Test_Mock_SageMakerExperiment_String(t *testing.T) {
	a := assert.New(t)
	experiment := SageMakerExperiment{ExperimentName: ptr.String("my-experiment")}
	a.Equal("my-experiment", experiment.String())
}
