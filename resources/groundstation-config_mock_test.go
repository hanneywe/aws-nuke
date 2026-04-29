package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/groundstation"
	groundstationtypes "github.com/aws/aws-sdk-go-v2/service/groundstation/types"
)

func Test_Mock_GroundStationConfig_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGroundStationClient)

	mockClient.On("ListConfigs", mock.Anything, mock.Anything).
		Return(&groundstation.ListConfigsOutput{
			ConfigList: []groundstationtypes.ConfigListItem{
				{
					ConfigId:   ptr.String("cfg-12345"),
					ConfigType: groundstationtypes.ConfigCapabilityTypeTracking,
					Name:       ptr.String("my-config"),
				},
			},
		}, nil)

	lister := &GroundStationConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGroundStationListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	cfg := resources[0].(*GroundStationConfig)
	a.Equal("cfg-12345", *cfg.ConfigID)
	a.Equal("my-config", *cfg.Name)
	a.Equal(groundstationtypes.ConfigCapabilityTypeTracking, cfg.ConfigType)

	mockClient.AssertExpectations(t)
}

func Test_Mock_GroundStationConfig_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGroundStationClient)

	mockClient.On("ListConfigs", mock.Anything, mock.Anything).
		Return(&groundstation.ListConfigsOutput{
			ConfigList: []groundstationtypes.ConfigListItem{},
		}, nil)

	lister := &GroundStationConfigLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGroundStationListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_GroundStationConfig_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGroundStationClient)

	cfg := &GroundStationConfig{
		svc:        mockClient,
		ConfigID:   ptr.String("cfg-12345"),
		ConfigType: groundstationtypes.ConfigCapabilityTypeTracking,
	}

	mockClient.On("DeleteConfig", mock.Anything, &groundstation.DeleteConfigInput{
		ConfigId:   cfg.ConfigID,
		ConfigType: cfg.ConfigType,
	}).Return(&groundstation.DeleteConfigOutput{}, nil)

	a.NoError(cfg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GroundStationConfig_Properties(t *testing.T) {
	a := assert.New(t)

	cfg := GroundStationConfig{
		ConfigID:   ptr.String("cfg-12345"),
		ConfigType: groundstationtypes.ConfigCapabilityTypeTracking,
		Name:       ptr.String("my-config"),
	}

	props := cfg.Properties()
	a.Equal("cfg-12345", props.Get("ConfigId"))
	a.Equal("my-config", props.Get("Name"))
}

func Test_Mock_GroundStationConfig_String(t *testing.T) {
	a := assert.New(t)
	cfg := GroundStationConfig{ConfigID: ptr.String("cfg-12345")}
	a.Equal("cfg-12345", cfg.String())
}
