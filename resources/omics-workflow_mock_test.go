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

func Test_Mock_OmicsWorkflow_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListWorkflows", mock.Anything, mock.Anything).
		Return(&omics.ListWorkflowsOutput{
			Items: []omicstypes.WorkflowListItem{
				{
					Id:   ptr.String("wf-1"),
					Name: ptr.String("my-workflow"),
				},
			},
		}, nil)

	lister := &OmicsWorkflowLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	workflow := resources[0].(*OmicsWorkflow)
	assertions.Equal("wf-1", *workflow.ID)
	assertions.Equal("my-workflow", *workflow.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsWorkflow_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	mockClient.
		On("ListWorkflows", mock.Anything, mock.Anything).
		Return(&omics.ListWorkflowsOutput{
			Items: []omicstypes.WorkflowListItem{},
		}, nil)

	lister := &OmicsWorkflowLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testOmicsListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsWorkflow_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockOmicsClient)

	workflow := &OmicsWorkflow{
		svc:  mockClient,
		ID:   ptr.String("wf-1"),
		Name: ptr.String("my-workflow"),
	}

	mockClient.
		On("DeleteWorkflow", mock.Anything, &omics.DeleteWorkflowInput{
			Id: workflow.ID,
		}).
		Return(&omics.DeleteWorkflowOutput{}, nil)

	err := workflow.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_OmicsWorkflow_Properties(t *testing.T) {
	assertions := assert.New(t)

	workflow := OmicsWorkflow{
		ID:   ptr.String("wf-1"),
		Name: ptr.String("my-workflow"),
	}

	properties := workflow.Properties()

	assertions.Equal("wf-1", properties.Get("ID"))
	assertions.Equal("my-workflow", properties.Get("Name"))
}

func Test_Mock_OmicsWorkflow_String(t *testing.T) {
	assertions := assert.New(t)

	workflow := OmicsWorkflow{
		ID: ptr.String("wf-1"),
	}

	assertions.Equal("wf-1", workflow.String())
}
