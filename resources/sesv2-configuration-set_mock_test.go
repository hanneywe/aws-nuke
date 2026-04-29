package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func Test_Mock_SESv2ConfigurationSet_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListConfigurationSets", mock.Anything, mock.Anything).
		Return(&sesv2.ListConfigurationSetsOutput{ConfigurationSets: []string{"my-config-set"}}, nil)
	lister := &SESv2ConfigurationSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	cs := resources[0].(*SESv2ConfigurationSet)
	a.Equal("my-config-set", *cs.ConfigurationSetName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2ConfigurationSet_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListConfigurationSets", mock.Anything, mock.Anything).
		Return(&sesv2.ListConfigurationSetsOutput{ConfigurationSets: []string{}}, nil)
	lister := &SESv2ConfigurationSetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2ConfigurationSet_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	cs := &SESv2ConfigurationSet{svc: mockClient, ConfigurationSetName: ptr.String("my-config-set")}
	mockClient.On("DeleteConfigurationSet", mock.Anything, &sesv2.DeleteConfigurationSetInput{ConfigurationSetName: cs.ConfigurationSetName}).
		Return(&sesv2.DeleteConfigurationSetOutput{}, nil)
	a.NoError(cs.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2ConfigurationSet_Properties(t *testing.T) {
	a := assert.New(t)
	cs := SESv2ConfigurationSet{ConfigurationSetName: ptr.String("my-config-set")}
	a.Equal("my-config-set", cs.Properties().Get("ConfigurationSetName"))
}

func Test_Mock_SESv2ConfigurationSet_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-config-set", (&SESv2ConfigurationSet{ConfigurationSetName: ptr.String("my-config-set")}).String())
}
