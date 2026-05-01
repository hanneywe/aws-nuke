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

func Test_Mock_SageMakerTrialComponent_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListTrialComponents", mock.Anything, mock.Anything).
		Return(&sagemaker.ListTrialComponentsOutput{
			TrialComponentSummaries: []sagemakertypes.TrialComponentSummary{
				{
					TrialComponentName: ptr.String("my-trial-component"),
					TrialComponentArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:experiment-trial-component/my-trial-component"),
				},
			},
		}, nil)

	lister := &SageMakerTrialComponentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	component := resources[0].(*SageMakerTrialComponent)
	a.Equal("my-trial-component", *component.TrialComponentName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:experiment-trial-component/my-trial-component", *component.TrialComponentArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerTrialComponent_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListTrialComponents", mock.Anything, mock.Anything).
		Return(&sagemaker.ListTrialComponentsOutput{
			TrialComponentSummaries: []sagemakertypes.TrialComponentSummary{},
		}, nil)

	lister := &SageMakerTrialComponentLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerTrialComponent_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	component := &SageMakerTrialComponent{
		svc:                mockClient,
		TrialComponentName: ptr.String("my-trial-component"),
	}

	mockClient.On("DeleteTrialComponent", mock.Anything, &sagemaker.DeleteTrialComponentInput{
		TrialComponentName: component.TrialComponentName,
	}).Return(&sagemaker.DeleteTrialComponentOutput{}, nil)

	a.NoError(component.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerTrialComponent_Properties(t *testing.T) {
	a := assert.New(t)

	component := SageMakerTrialComponent{
		TrialComponentName: ptr.String("my-trial-component"),
		TrialComponentArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:experiment-trial-component/my-trial-component"),
	}

	props := component.Properties()
	a.Equal("my-trial-component", props.Get("TrialComponentName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:experiment-trial-component/my-trial-component", props.Get("TrialComponentArn"))
}

func Test_Mock_SageMakerTrialComponent_String(t *testing.T) {
	a := assert.New(t)
	component := SageMakerTrialComponent{TrialComponentName: ptr.String("my-trial-component")}
	a.Equal("my-trial-component", component.String())
}
