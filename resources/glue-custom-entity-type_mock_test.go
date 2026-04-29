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

func Test_Mock_GlueCustomEntityType_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListCustomEntityTypes", mock.Anything, mock.Anything).
		Return(&glue.ListCustomEntityTypesOutput{
			CustomEntityTypes: []gluetypes.CustomEntityType{
				{Name: ptr.String("my-entity-type")},
			},
		}, nil)
	lister := &GlueCustomEntityTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-entity-type", resources[0].(*GlueCustomEntityType).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueCustomEntityType_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListCustomEntityTypes", mock.Anything, mock.Anything).
		Return(&glue.ListCustomEntityTypesOutput{CustomEntityTypes: []gluetypes.CustomEntityType{}}, nil)
	lister := &GlueCustomEntityTypeLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueCustomEntityType_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	r := &GlueCustomEntityType{svc: mockClient, Name: ptr.String("my-entity-type")}
	mockClient.On("DeleteCustomEntityType", mock.Anything, &glue.DeleteCustomEntityTypeInput{
		Name: r.Name,
	}).Return(&glue.DeleteCustomEntityTypeOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueCustomEntityType_Properties(t *testing.T) {
	a := assert.New(t)
	r := GlueCustomEntityType{Name: ptr.String("my-entity-type")}
	a.Equal("my-entity-type", r.Properties().Get("Name"))
}

func Test_Mock_GlueCustomEntityType_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-entity-type", (&GlueCustomEntityType{Name: ptr.String("my-entity-type")}).String())
}
