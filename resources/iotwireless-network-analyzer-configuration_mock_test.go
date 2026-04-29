package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iotwireless"
	iotwirelesstypes "github.com/aws/aws-sdk-go-v2/service/iotwireless/types"
)

func Test_Mock_IoTWirelessNetworkAnalyzerConfiguration_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListNetworkAnalyzerConfigurations", mock.Anything, mock.Anything).
		Return(&iotwireless.ListNetworkAnalyzerConfigurationsOutput{
			NetworkAnalyzerConfigurationList: []iotwirelesstypes.NetworkAnalyzerConfigurations{
				{
					Name: ptr.String("my-config"),
				},
			},
		}, nil)

	lister := &IoTWirelessNetworkAnalyzerConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cfg := resources[0].(*IoTWirelessNetworkAnalyzerConfiguration)
	a.Equal("my-config", *cfg.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessNetworkAnalyzerConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	mockClient.On("ListNetworkAnalyzerConfigurations", mock.Anything, mock.Anything).
		Return(&iotwireless.ListNetworkAnalyzerConfigurationsOutput{
			NetworkAnalyzerConfigurationList: []iotwirelesstypes.NetworkAnalyzerConfigurations{},
		}, nil)

	lister := &IoTWirelessNetworkAnalyzerConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTWirelessListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessNetworkAnalyzerConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIoTWirelessClient)

	cfg := &IoTWirelessNetworkAnalyzerConfiguration{
		svc:  mockClient,
		Name: ptr.String("my-config"),
	}

	mockClient.On("DeleteNetworkAnalyzerConfiguration", mock.Anything,
		&iotwireless.DeleteNetworkAnalyzerConfigurationInput{
			ConfigurationName: cfg.Name,
		}).Return(&iotwireless.DeleteNetworkAnalyzerConfigurationOutput{}, nil)

	a.NoError(cfg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTWirelessNetworkAnalyzerConfiguration_Properties(t *testing.T) {
	a := assert.New(t)

	cfg := IoTWirelessNetworkAnalyzerConfiguration{
		Name: ptr.String("my-config"),
	}

	props := cfg.Properties()
	a.Equal("my-config", props.Get("Name"))
}

func Test_Mock_IoTWirelessNetworkAnalyzerConfiguration_String(t *testing.T) {
	a := assert.New(t)
	cfg := IoTWirelessNetworkAnalyzerConfiguration{Name: ptr.String("my-config")}
	a.Equal("my-config", cfg.String())
}
