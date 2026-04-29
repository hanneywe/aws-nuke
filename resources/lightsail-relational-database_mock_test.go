package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lightsailtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

func Test_Mock_LightsailRelationalDatabase_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetRelationalDatabases", mock.Anything, mock.Anything).
		Return(&lightsail.GetRelationalDatabasesOutput{
			RelationalDatabases: []lightsailtypes.RelationalDatabase{
				{Name: ptr.String("test-value"), Engine: ptr.String("test-value"), EngineVersion: ptr.String("test-value")},
			},
		}, nil)

	lister := &LightsailRelationalDatabaseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*LightsailRelationalDatabase)
	a.Equal("test-value", *r.RelationalDatabaseName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailRelationalDatabase_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	mockClient.On("GetRelationalDatabases", mock.Anything, mock.Anything).
		Return(&lightsail.GetRelationalDatabasesOutput{
			RelationalDatabases: []lightsailtypes.RelationalDatabase{},
		}, nil)

	lister := &LightsailRelationalDatabaseLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLightsailListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailRelationalDatabase_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLightsailClient)

	r := &LightsailRelationalDatabase{
		svc:                    mockClient,
		RelationalDatabaseName: ptr.String("test-relationaldatabasename"),
	}

	mockClient.On("DeleteRelationalDatabase", mock.Anything,
		&lightsail.DeleteRelationalDatabaseInput{
			RelationalDatabaseName: r.RelationalDatabaseName,
			SkipFinalSnapshot:      aws.Bool(true),
		}).Return(&lightsail.DeleteRelationalDatabaseOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LightsailRelationalDatabase_Properties(t *testing.T) {
	a := assert.New(t)
	r := &LightsailRelationalDatabase{
		RelationalDatabaseName: ptr.String("test-relationaldatabasename"),
		Engine:                 ptr.String("test-engine"),
		EngineVersion:          ptr.String("test-engineversion"),
	}
	props := r.Properties()
	a.Equal("test-relationaldatabasename", props.Get("RelationalDatabaseName"))
	a.Equal("test-engine", props.Get("Engine"))
	a.Equal("test-engineversion", props.Get("EngineVersion"))
}

func Test_Mock_LightsailRelationalDatabase_String(t *testing.T) {
	a := assert.New(t)
	r := &LightsailRelationalDatabase{
		RelationalDatabaseName: ptr.String("test-relationaldatabasename"),
	}
	a.Equal("test-relationaldatabasename", r.String())
}
