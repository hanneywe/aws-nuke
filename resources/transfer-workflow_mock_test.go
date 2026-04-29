package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/transfer"
	transfertypes "github.com/aws/aws-sdk-go-v2/service/transfer/types"
)

func Test_Mock_TransferWorkflow_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.On("ListWorkflows", mock.Anything, mock.Anything).
		Return(&transfer.ListWorkflowsOutput{
			Workflows: []transfertypes.ListedWorkflow{
				{WorkflowId: ptr.String("test-value"), Description: ptr.String("test-value"), Arn: ptr.String("test-value")},
			},
		}, nil)

	lister := &TransferWorkflowLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*TransferWorkflow)
	a.Equal("test-value", *r.WorkflowID)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferWorkflow_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	mockClient.On("ListWorkflows", mock.Anything, mock.Anything).
		Return(&transfer.ListWorkflowsOutput{
			Workflows: []transfertypes.ListedWorkflow{},
		}, nil)

	lister := &TransferWorkflowLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testTransferListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferWorkflow_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockTransferClient)

	r := &TransferWorkflow{
		svc:        mockClient,
		WorkflowID: ptr.String("test-workflowid"),
	}

	mockClient.On("DeleteWorkflow", mock.Anything,
		&transfer.DeleteWorkflowInput{
			WorkflowId: r.WorkflowID,
		}).Return(&transfer.DeleteWorkflowOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_TransferWorkflow_Properties(t *testing.T) {
	a := assert.New(t)
	r := &TransferWorkflow{
		WorkflowID:  ptr.String("test-workflowid"),
		Description: ptr.String("test-description"),
		Arn:         ptr.String("test-arn"),
	}
	props := r.Properties()
	a.Equal("test-workflowid", props.Get("WorkflowID"))
	a.Equal("test-description", props.Get("Description"))
	a.Equal("test-arn", props.Get("Arn"))
}

func Test_Mock_TransferWorkflow_String(t *testing.T) {
	a := assert.New(t)
	r := &TransferWorkflow{
		WorkflowID: ptr.String("test-workflowid"),
	}
	a.Equal("test-workflowid", r.String())
}
