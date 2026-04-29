package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/schemas"
	schemastypes "github.com/aws/aws-sdk-go-v2/service/schemas/types"
)

func Test_Mock_SchemasRegistry_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSchemasClient)
	mockClient.On("ListRegistries", mock.Anything, mock.Anything).
		Return(&schemas.ListRegistriesOutput{
			Registries: []schemastypes.RegistrySummary{
				{RegistryName: ptr.String("my-registry"), RegistryArn: ptr.String("arn:aws:schemas:us-east-1:123456789012:registry/my-registry")},
			},
		}, nil)
	lister := &SchemasRegistryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSchemasListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	reg := resources[0].(*SchemasRegistry)
	a.Equal("my-registry", *reg.RegistryName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchemasRegistry_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSchemasClient)
	mockClient.On("ListRegistries", mock.Anything, mock.Anything).
		Return(&schemas.ListRegistriesOutput{Registries: []schemastypes.RegistrySummary{}}, nil)
	lister := &SchemasRegistryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSchemasListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchemasRegistry_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSchemasClient)
	reg := &SchemasRegistry{svc: mockClient, RegistryName: ptr.String("my-registry")}
	mockClient.On("DeleteRegistry", mock.Anything, &schemas.DeleteRegistryInput{RegistryName: reg.RegistryName}).
		Return(&schemas.DeleteRegistryOutput{}, nil)
	a.NoError(reg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SchemasRegistry_Filter_AWSManaged(t *testing.T) {
	a := assert.New(t)
	a.Error((&SchemasRegistry{RegistryName: ptr.String("aws.events")}).Filter())
	a.Error((&SchemasRegistry{RegistryName: ptr.String("discovered-schemas")}).Filter())
}

func Test_Mock_SchemasRegistry_Filter_Custom(t *testing.T) {
	a := assert.New(t)
	a.NoError((&SchemasRegistry{RegistryName: ptr.String("my-registry")}).Filter())
}

func Test_Mock_SchemasRegistry_Properties(t *testing.T) {
	a := assert.New(t)
	reg := SchemasRegistry{
		RegistryName: ptr.String("my-registry"),
		RegistryArn:  ptr.String("arn:aws:schemas:us-east-1:123456789012:registry/my-registry"),
	}
	a.Equal("my-registry", reg.Properties().Get("RegistryName"))
}

func Test_Mock_SchemasRegistry_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-registry", (&SchemasRegistry{RegistryName: ptr.String("my-registry")}).String())
}
