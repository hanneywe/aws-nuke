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

func Test_Mock_AppRunnerObservabilityConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRunnerClient)
	mockClient.On("ListObservabilityConfigurations", mock.Anything, mock.Anything).
		Return(&apprunner.ListObservabilityConfigurationsOutput{
			ObservabilityConfigurationSummaryList: []apprunnertypes.ObservabilityConfigurationSummary{
				{
					ObservabilityConfigurationArn:  ptr.String("arn:aws:apprunner:us-east-1:123456789012:observabilityconfiguration/my-oc"),
					ObservabilityConfigurationName: ptr.String("my-oc"),
				},
			},
		}, nil)
	lister := &AppRunnerObservabilityConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppRunnerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	oc := resources[0].(*AppRunnerObservabilityConfiguration)
	a.Equal("my-oc", *oc.ObservabilityConfigurationName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRunnerObservabilityConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRunnerClient)
	mockClient.On("ListObservabilityConfigurations", mock.Anything, mock.Anything).
		Return(&apprunner.ListObservabilityConfigurationsOutput{
			ObservabilityConfigurationSummaryList: []apprunnertypes.ObservabilityConfigurationSummary{},
		}, nil)
	lister := &AppRunnerObservabilityConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAppRunnerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRunnerObservabilityConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAppRunnerClient)
	ocArn := ptr.String("arn:aws:apprunner:us-east-1:123456789012:observabilityconfiguration/my-oc")
	oc := &AppRunnerObservabilityConfiguration{
		svc:                           mockClient,
		ObservabilityConfigurationArn: ocArn,
	}
	mockClient.On("DeleteObservabilityConfiguration", mock.Anything,
		&apprunner.DeleteObservabilityConfigurationInput{
			ObservabilityConfigurationArn: oc.ObservabilityConfigurationArn,
		}).
		Return(&apprunner.DeleteObservabilityConfigurationOutput{}, nil)
	a.NoError(oc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AppRunnerObservabilityConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	oc := AppRunnerObservabilityConfiguration{
		ObservabilityConfigurationArn:  ptr.String("arn"),
		ObservabilityConfigurationName: ptr.String("my-oc"),
	}
	a.Equal("my-oc", oc.Properties().Get("ObservabilityConfigurationName"))
}

func Test_Mock_AppRunnerObservabilityConfiguration_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-oc", (&AppRunnerObservabilityConfiguration{ObservabilityConfigurationName: ptr.String("my-oc")}).String())
}

func Test_Mock_AppRunnerObservabilityConfiguration_Filter_Default(t *testing.T) {
	a := assert.New(t)
	r := AppRunnerObservabilityConfiguration{ObservabilityConfigurationName: ptr.String("DefaultConfiguration")}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "cannot delete default configuration")
}

func Test_Mock_AppRunnerObservabilityConfiguration_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	r := AppRunnerObservabilityConfiguration{ObservabilityConfigurationName: ptr.String("my-oc")}
	a.NoError(r.Filter())
}
