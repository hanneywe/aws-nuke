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

func Test_Mock_GlueUserDefinedFunction_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("GetDatabases", mock.Anything, mock.Anything).
		Return(&glue.GetDatabasesOutput{
			DatabaseList: []gluetypes.Database{
				{Name: ptr.String("mydb")},
			},
		}, nil)
	mockClient.On("GetUserDefinedFunctions", mock.Anything, mock.Anything).
		Return(&glue.GetUserDefinedFunctionsOutput{
			UserDefinedFunctions: []gluetypes.UserDefinedFunction{
				{
					DatabaseName: ptr.String("mydb"),
					FunctionName: ptr.String("myfunc"),
					CatalogId:    ptr.String("123456789012"),
				},
			},
		}, nil)
	lister := &GlueUserDefinedFunctionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("mydb/myfunc", resources[0].(*GlueUserDefinedFunction).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUserDefinedFunction_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("GetDatabases", mock.Anything, mock.Anything).
		Return(&glue.GetDatabasesOutput{DatabaseList: []gluetypes.Database{}}, nil)
	lister := &GlueUserDefinedFunctionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUserDefinedFunction_List_MultipleDatabases(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("GetDatabases", mock.Anything, mock.Anything).
		Return(&glue.GetDatabasesOutput{
			DatabaseList: []gluetypes.Database{
				{Name: ptr.String("db1")},
				{Name: ptr.String("db2")},
			},
		}, nil)
	mockClient.On("GetUserDefinedFunctions", mock.Anything, mock.MatchedBy(func(input *glue.GetUserDefinedFunctionsInput) bool {
		return *input.DatabaseName == "db1"
	})).Return(&glue.GetUserDefinedFunctionsOutput{
		UserDefinedFunctions: []gluetypes.UserDefinedFunction{
			{DatabaseName: ptr.String("db1"), FunctionName: ptr.String("func1"), CatalogId: ptr.String("123456789012")},
		},
	}, nil)
	mockClient.On("GetUserDefinedFunctions", mock.Anything, mock.MatchedBy(func(input *glue.GetUserDefinedFunctionsInput) bool {
		return *input.DatabaseName == "db2"
	})).Return(&glue.GetUserDefinedFunctionsOutput{
		UserDefinedFunctions: []gluetypes.UserDefinedFunction{
			{DatabaseName: ptr.String("db2"), FunctionName: ptr.String("func2"), CatalogId: ptr.String("123456789012")},
		},
	}, nil)
	lister := &GlueUserDefinedFunctionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 2)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUserDefinedFunction_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	r := &GlueUserDefinedFunction{
		svc:          mockClient,
		DatabaseName: ptr.String("mydb"),
		FunctionName: ptr.String("myfunc"),
	}
	mockClient.On("DeleteUserDefinedFunction", mock.Anything, &glue.DeleteUserDefinedFunctionInput{
		DatabaseName: r.DatabaseName,
		FunctionName: r.FunctionName,
	}).Return(&glue.DeleteUserDefinedFunctionOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueUserDefinedFunction_Properties(t *testing.T) {
	a := assert.New(t)
	r := GlueUserDefinedFunction{
		DatabaseName: ptr.String("mydb"),
		FunctionName: ptr.String("myfunc"),
		CatalogID:    ptr.String("123456789012"),
	}
	a.Equal("mydb", r.Properties().Get("DatabaseName"))
	a.Equal("myfunc", r.Properties().Get("FunctionName"))
	a.Equal("123456789012", r.Properties().Get("CatalogId"))
}

func Test_Mock_GlueUserDefinedFunction_String(t *testing.T) {
	a := assert.New(t)
	r := &GlueUserDefinedFunction{
		DatabaseName: ptr.String("mydb"),
		FunctionName: ptr.String("myfunc"),
	}
	a.Equal("mydb/myfunc", r.String())
}
