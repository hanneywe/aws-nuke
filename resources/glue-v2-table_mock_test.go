package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
)

func Test_Mock_GlueTable_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	mockClient.On("GetDatabases", mock.Anything, mock.Anything).
		Return(&glue.GetDatabasesOutput{
			DatabaseList: []gluetypes.Database{
				{Name: ptr.String("test-db")},
			},
		}, nil)

	mockClient.On("GetTables", mock.Anything, mock.Anything).
		Return(&glue.GetTablesOutput{
			TableList: []gluetypes.Table{
				{Name: ptr.String("test-table")},
			},
		}, nil)

	lister := &GlueTableLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*GlueTable)
	a.Equal("test-db", *r.DatabaseName)
	a.Equal("test-table", *r.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueTable_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	mockClient.On("GetDatabases", mock.Anything, mock.Anything).
		Return(&glue.GetDatabasesOutput{
			DatabaseList: []gluetypes.Database{},
		}, nil)

	lister := &GlueTableLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueTable_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)

	r := &GlueTable{
		svc:          mockClient,
		DatabaseName: ptr.String("test-db"),
		Name:         ptr.String("test-table"),
	}

	mockClient.On("DeleteTable", mock.Anything,
		&glue.DeleteTableInput{
			DatabaseName: r.DatabaseName,
			Name:         r.Name,
		}).Return(&glue.DeleteTableOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueTable_Properties(t *testing.T) {
	a := assert.New(t)
	r := &GlueTable{
		DatabaseName: ptr.String("test-db"),
		Name:         ptr.String("test-table"),
	}
	props := r.Properties()
	a.Equal("test-db", props.Get("DatabaseName"))
	a.Equal("test-table", props.Get("Name"))
}

func Test_Mock_GlueTable_String(t *testing.T) {
	a := assert.New(t)
	r := &GlueTable{
		DatabaseName: ptr.String("test-db"),
		Name:         ptr.String("test-table"),
	}
	a.Equal(fmt.Sprintf("%s/%s", "test-db", "test-table"), r.String())
}
