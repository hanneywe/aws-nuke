package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/entityresolution"
	entityresolutiontypes "github.com/aws/aws-sdk-go-v2/service/entityresolution/types"
)

func Test_Mock_EntityResolutionSchemaMapping_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEntityresolutionClient)

	mockClient.On("ListSchemaMappings", mock.Anything, mock.Anything).
		Return(&entityresolution.ListSchemaMappingsOutput{
			SchemaList: []entityresolutiontypes.SchemaMappingSummary{
				{SchemaName: ptr.String("test-value"), SchemaArn: ptr.String("test-value")},
			},
		}, nil)

	lister := &EntityResolutionSchemaMappingLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEntityresolutionListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*EntityResolutionSchemaMapping)
	a.Equal("test-value", *r.SchemaName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EntityResolutionSchemaMapping_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEntityresolutionClient)

	mockClient.On("ListSchemaMappings", mock.Anything, mock.Anything).
		Return(&entityresolution.ListSchemaMappingsOutput{
			SchemaList: []entityresolutiontypes.SchemaMappingSummary{},
		}, nil)

	lister := &EntityResolutionSchemaMappingLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEntityresolutionListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_EntityResolutionSchemaMapping_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEntityresolutionClient)

	r := &EntityResolutionSchemaMapping{
		svc:        mockClient,
		SchemaName: ptr.String("test-schemaname"),
	}

	mockClient.On("DeleteSchemaMapping", mock.Anything,
		&entityresolution.DeleteSchemaMappingInput{
			SchemaName: r.SchemaName,
		}).Return(&entityresolution.DeleteSchemaMappingOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_EntityResolutionSchemaMapping_Properties(t *testing.T) {
	a := assert.New(t)
	r := &EntityResolutionSchemaMapping{
		SchemaName: ptr.String("test-schemaname"),
		SchemaArn:  ptr.String("test-schemaarn"),
	}
	props := r.Properties()
	a.Equal("test-schemaname", props.Get("SchemaName"))
	a.Equal("test-schemaarn", props.Get("SchemaArn"))
}

func Test_Mock_EntityResolutionSchemaMapping_String(t *testing.T) {
	a := assert.New(t)
	r := &EntityResolutionSchemaMapping{
		SchemaName: ptr.String("test-schemaname"),
	}
	a.Equal("test-schemaname", r.String())
}
