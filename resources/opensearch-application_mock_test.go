package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	opensearchtypes "github.com/aws/aws-sdk-go-v2/service/opensearch/types"
)

func Test_Mock_OpenSearchApplication_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOpenSearchClient)

	mockClient.On("ListApplications", mock.Anything, mock.Anything).
		Return(&opensearch.ListApplicationsOutput{
			ApplicationSummaries: []opensearchtypes.ApplicationSummary{
				{
					Id:   ptr.String("app-abc123"),
					Name: ptr.String("my-app"),
				},
			},
		}, nil)

	lister := &OpenSearchApplicationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOpenSearchListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	app := resources[0].(*OpenSearchApplication)
	a.Equal("app-abc123", *app.ID)
	a.Equal("my-app", *app.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OpenSearchApplication_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOpenSearchClient)

	mockClient.On("ListApplications", mock.Anything, mock.Anything).
		Return(&opensearch.ListApplicationsOutput{
			ApplicationSummaries: []opensearchtypes.ApplicationSummary{},
		}, nil)

	lister := &OpenSearchApplicationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testOpenSearchListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OpenSearchApplication_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockOpenSearchClient)

	app := &OpenSearchApplication{
		svc: mockClient,
		ID:  ptr.String("app-abc123"),
	}

	mockClient.On("DeleteApplication", mock.Anything, &opensearch.DeleteApplicationInput{
		Id: app.ID,
	}).Return(&opensearch.DeleteApplicationOutput{}, nil)

	a.NoError(app.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_OpenSearchApplication_Properties(t *testing.T) {
	a := assert.New(t)

	app := OpenSearchApplication{
		ID:   ptr.String("app-abc123"),
		Name: ptr.String("my-app"),
	}

	props := app.Properties()
	a.Equal("app-abc123", props.Get("ID"))
	a.Equal("my-app", props.Get("Name"))
}

func Test_Mock_OpenSearchApplication_String(t *testing.T) {
	a := assert.New(t)
	app := OpenSearchApplication{ID: ptr.String("app-abc123")}
	a.Equal("app-abc123", app.String())
}
