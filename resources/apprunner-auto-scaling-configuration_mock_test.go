package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	apprunnertypes "github.com/aws/aws-sdk-go-v2/service/apprunner/types"
)

func Test_Mock_AppRunnerAutoScalingConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRunnerClient)
	mockClient.On("ListAutoScalingConfigurations", mock.Anything, mock.Anything).
		Return(&apprunner.ListAutoScalingConfigurationsOutput{
			AutoScalingConfigurationSummaryList: []apprunnertypes.AutoScalingConfigurationSummary{
				{
					AutoScalingConfigurationArn:  ptr.String("arn:aws:apprunner:us-east-1:123456789012:autoscalingconfiguration/my-asc"),
					AutoScalingConfigurationName: ptr.String("my-asc"),
				},
			},
		}, nil)
	lister := &AppRunnerAutoScalingConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppRunnerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	asc := resources[0].(*AppRunnerAutoScalingConfiguration)
	a.Equal("my-asc", *asc.AutoScalingConfigurationName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRunnerAutoScalingConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRunnerClient)
	mockClient.On("ListAutoScalingConfigurations", mock.Anything, mock.Anything).
		Return(&apprunner.ListAutoScalingConfigurationsOutput{
			AutoScalingConfigurationSummaryList: []apprunnertypes.AutoScalingConfigurationSummary{},
		}, nil)
	lister := &AppRunnerAutoScalingConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppRunnerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRunnerAutoScalingConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRunnerClient)
	ascArn := ptr.String("arn:aws:apprunner:us-east-1:123456789012:autoscalingconfiguration/my-asc")
	asc := &AppRunnerAutoScalingConfiguration{
		svc:                         mockClient,
		AutoScalingConfigurationArn: ascArn,
	}
	mockClient.On("DeleteAutoScalingConfiguration", mock.Anything,
		&apprunner.DeleteAutoScalingConfigurationInput{
			AutoScalingConfigurationArn: asc.AutoScalingConfigurationArn,
		}).
		Return(&apprunner.DeleteAutoScalingConfigurationOutput{}, nil)
	a.NoError(asc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRunnerAutoScalingConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	asc := AppRunnerAutoScalingConfiguration{
		AutoScalingConfigurationArn:  ptr.String("arn"),
		AutoScalingConfigurationName: ptr.String("my-asc"),
	}
	a.Equal("my-asc", asc.Properties().Get("AutoScalingConfigurationName"))
}

func Test_Mock_AppRunnerAutoScalingConfiguration_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-asc", (&AppRunnerAutoScalingConfiguration{AutoScalingConfigurationName: ptr.String("my-asc")}).String())
}

func Test_Mock_AppRunnerAutoScalingConfiguration_Filter_Default(t *testing.T) {
	a := assert.New(t)
	r := AppRunnerAutoScalingConfiguration{AutoScalingConfigurationName: ptr.String("DefaultConfiguration")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete default configuration")
}

func Test_Mock_AppRunnerAutoScalingConfiguration_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := AppRunnerAutoScalingConfiguration{AutoScalingConfigurationName: ptr.String("my-asc")}
	a.NoError(r.Filter())
}
