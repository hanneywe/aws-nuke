package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
)

func Test_Mock_GlueRegistry_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListRegistries", mock.Anything, mock.Anything).
		Return(&glue.ListRegistriesOutput{
			Registries: []gluetypes.RegistryListItem{
				{RegistryName: ptr.String("my-registry"), RegistryArn: ptr.String("arn:aws:glue:us-east-1:123456789012:registry/my-registry")},
			},
		}, nil)
	lister := &GlueRegistryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-registry", resources[0].(*GlueRegistry).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueRegistry_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListRegistries", mock.Anything, mock.Anything).
		Return(&glue.ListRegistriesOutput{Registries: []gluetypes.RegistryListItem{}}, nil)
	lister := &GlueRegistryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueRegistry_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	r := &GlueRegistry{
		svc:          mockClient,
		RegistryName: ptr.String("my-registry"),
		RegistryArn:  ptr.String("arn:aws:glue:us-east-1:123456789012:registry/my-registry"),
	}
	mockClient.On("DeleteRegistry", mock.Anything, &glue.DeleteRegistryInput{
		RegistryId: &gluetypes.RegistryId{RegistryArn: r.RegistryArn},
	}).Return(&glue.DeleteRegistryOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueRegistry_Properties(t *testing.T) {
	a := assert.New(t)
	r := GlueRegistry{
		RegistryName: ptr.String("my-registry"),
		RegistryArn:  ptr.String("arn:aws:glue:us-east-1:123456789012:registry/my-registry"),
	}
	a.Equal("my-registry", r.Properties().Get("RegistryName"))
	a.Equal("arn:aws:glue:us-east-1:123456789012:registry/my-registry", r.Properties().Get("RegistryArn"))
}

func Test_Mock_GlueRegistry_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-registry", (&GlueRegistry{RegistryName: ptr.String("my-registry")}).String())
}
