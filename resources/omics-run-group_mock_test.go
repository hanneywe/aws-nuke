package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/omics"
	omicstypes "github.com/aws/aws-sdk-go-v2/service/omics/types"
)

func Test_Mock_OmicsRunGroup_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListRunGroups", mock.Anything, mock.Anything).
		Return(&omics.ListRunGroupsOutput{
			Items: []omicstypes.RunGroupListItem{
				{
					Id:   ptr.String("rg-1"),
					Name: ptr.String("my-run-group"),
				},
			},
		}, nil)

	lister := &OmicsRunGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	runGroup := resources[0].(*OmicsRunGroup)
	assertions.Equal("rg-1", *runGroup.ID)
	assertions.Equal("my-run-group", *runGroup.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsRunGroup_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListRunGroups", mock.Anything, mock.Anything).
		Return(&omics.ListRunGroupsOutput{
			Items: []omicstypes.RunGroupListItem{},
		}, nil)

	lister := &OmicsRunGroupLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsRunGroup_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	runGroup := &OmicsRunGroup{
		svc:  mockClient,
		ID:   ptr.String("rg-1"),
		Name: ptr.String("my-run-group"),
	}

	mockClient.
		On("DeleteRunGroup", mock.Anything, &omics.DeleteRunGroupInput{
			Id: runGroup.ID,
		}).
		Return(&omics.DeleteRunGroupOutput{}, nil)

	err := runGroup.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsRunGroup_Properties(t *testing.T) {
	assertions := assert.New(t)

	runGroup := OmicsRunGroup{
		ID:   ptr.String("rg-1"),
		Name: ptr.String("my-run-group"),
	}

	properties := runGroup.Properties()

	assertions.Equal("rg-1", properties.Get("ID"))
	assertions.Equal("my-run-group", properties.Get("Name"))
}

func Test_Mock_OmicsRunGroup_String(t *testing.T) {
	assertions := assert.New(t)

	runGroup := OmicsRunGroup{
		ID: ptr.String("rg-1"),
	}

	assertions.Equal("rg-1", runGroup.String())
}
