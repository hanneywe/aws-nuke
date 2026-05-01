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

func Test_Mock_SageMakerHumanTaskUi_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListHumanTaskUis", mock.Anything, mock.Anything).
		Return(&sagemaker.ListHumanTaskUisOutput{
			HumanTaskUiSummaries: []sagemakertypes.HumanTaskUiSummary{
				{
					HumanTaskUiName: ptr.String("my-task-ui"),
					HumanTaskUiArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:human-task-ui/my-task-ui"),
				},
			},
		}, nil)

	lister := &SageMakerHumanTaskUILister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	ui := resources[0].(*SageMakerHumanTaskUI)
	a.Equal("my-task-ui", *ui.HumanTaskUIName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:human-task-ui/my-task-ui", *ui.HumanTaskUIArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerHumanTaskUi_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListHumanTaskUis", mock.Anything, mock.Anything).
		Return(&sagemaker.ListHumanTaskUisOutput{
			HumanTaskUiSummaries: []sagemakertypes.HumanTaskUiSummary{},
		}, nil)

	lister := &SageMakerHumanTaskUILister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerHumanTaskUi_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	ui := &SageMakerHumanTaskUI{
		svc:             mockClient,
		HumanTaskUIName: ptr.String("my-task-ui"),
	}

	mockClient.On("DeleteHumanTaskUi", mock.Anything, &sagemaker.DeleteHumanTaskUiInput{
		HumanTaskUiName: ui.HumanTaskUIName,
	}).Return(&sagemaker.DeleteHumanTaskUiOutput{}, nil)

	a.NoError(ui.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerHumanTaskUi_Properties(t *testing.T) {
	a := assert.New(t)

	ui := SageMakerHumanTaskUI{
		HumanTaskUIName: ptr.String("my-task-ui"),
		HumanTaskUIArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:human-task-ui/my-task-ui"),
	}

	props := ui.Properties()
	a.Equal("my-task-ui", props.Get("HumanTaskUIName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:human-task-ui/my-task-ui", props.Get("HumanTaskUIArn"))
}

func Test_Mock_SageMakerHumanTaskUi_String(t *testing.T) {
	a := assert.New(t)
	ui := SageMakerHumanTaskUI{HumanTaskUIName: ptr.String("my-task-ui")}
	a.Equal("my-task-ui", ui.String())
}
