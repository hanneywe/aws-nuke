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

func Test_Mock_SageMakerStudioLifecycleConfig_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListStudioLifecycleConfigs", mock.Anything, mock.Anything).
		Return(&sagemaker.ListStudioLifecycleConfigsOutput{
			StudioLifecycleConfigs: []sagemakertypes.StudioLifecycleConfigDetails{
				{
					StudioLifecycleConfigName: ptr.String("my-lc-config"),
					StudioLifecycleConfigArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:studio-lifecycle-config/my-lc-config"),
				},
			},
		}, nil)

	lister := &SageMakerStudioLifecycleConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	config := resources[0].(*SageMakerStudioLifecycleConfig)
	a.Equal("my-lc-config", *config.StudioLifecycleConfigName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:studio-lifecycle-config/my-lc-config", *config.StudioLifecycleConfigArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerStudioLifecycleConfig_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListStudioLifecycleConfigs", mock.Anything, mock.Anything).
		Return(&sagemaker.ListStudioLifecycleConfigsOutput{
			StudioLifecycleConfigs: []sagemakertypes.StudioLifecycleConfigDetails{},
		}, nil)

	lister := &SageMakerStudioLifecycleConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerStudioLifecycleConfig_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	config := &SageMakerStudioLifecycleConfig{
		svc:                       mockClient,
		StudioLifecycleConfigName: ptr.String("my-lc-config"),
	}

	mockClient.On("DeleteStudioLifecycleConfig", mock.Anything, &sagemaker.DeleteStudioLifecycleConfigInput{
		StudioLifecycleConfigName: config.StudioLifecycleConfigName,
	}).Return(&sagemaker.DeleteStudioLifecycleConfigOutput{}, nil)

	a.NoError(config.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerStudioLifecycleConfig_Properties(t *testing.T) {
	a := assert.New(t)

	config := SageMakerStudioLifecycleConfig{
		StudioLifecycleConfigName: ptr.String("my-lc-config"),
		StudioLifecycleConfigArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:studio-lifecycle-config/my-lc-config"),
	}

	props := config.Properties()
	a.Equal("my-lc-config", props.Get("StudioLifecycleConfigName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:studio-lifecycle-config/my-lc-config", props.Get("StudioLifecycleConfigArn"))
}

func Test_Mock_SageMakerStudioLifecycleConfig_String(t *testing.T) {
	a := assert.New(t)
	config := SageMakerStudioLifecycleConfig{StudioLifecycleConfigName: ptr.String("my-lc-config")}
	a.Equal("my-lc-config", config.String())
}
