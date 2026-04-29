package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/applicationdiscoveryservice"
)

func Test_Mock_DiscoveryApplication_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApplicationdiscoveryserviceClient)

	mockClient.On("ListConfigurations", mock.Anything, mock.Anything).
		Return(&applicationdiscoveryservice.ListConfigurationsOutput{
			Configurations: []map[string]string{
				{
					"application.configurationId": "test-id",
					"application.name":            "test-app",
				},
			},
		}, nil)

	lister := &DiscoveryApplicationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testApplicationdiscoveryserviceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*DiscoveryApplication)
	a.Equal("test-id", *r.ConfigurationID)
	a.Equal("test-app", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DiscoveryApplication_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApplicationdiscoveryserviceClient)

	mockClient.On("ListConfigurations", mock.Anything, mock.Anything).
		Return(&applicationdiscoveryservice.ListConfigurationsOutput{
			Configurations: []map[string]string{},
		}, nil)

	lister := &DiscoveryApplicationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testApplicationdiscoveryserviceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_DiscoveryApplication_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockApplicationdiscoveryserviceClient)

	r := &DiscoveryApplication{
		svc:             mockClient,
		ConfigurationID: ptr.String("test-configurationid"),
	}

	mockClient.On("DeleteApplications", mock.Anything,
		&applicationdiscoveryservice.DeleteApplicationsInput{
			ConfigurationIds: []string{"test-configurationid"},
		}).Return(&applicationdiscoveryservice.DeleteApplicationsOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_DiscoveryApplication_Properties(t *testing.T) {
	a := assert.New(t)
	r := &DiscoveryApplication{
		ConfigurationID: ptr.String("test-configurationid"),
		Name:            ptr.String("test-name"),
	}
	props := r.Properties()
	a.Equal("test-configurationid", props.Get("ConfigurationId"))
	a.Equal("test-name", props.Get("Name"))
}

func Test_Mock_DiscoveryApplication_String(t *testing.T) {
	a := assert.New(t)
	r := &DiscoveryApplication{
		ConfigurationID: ptr.String("test-configurationid"),
	}
	a.Equal("test-configurationid", r.String())
}
