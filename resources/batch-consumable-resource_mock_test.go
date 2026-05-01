package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"
)

func Test_Mock_BatchConsumableResource_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)
	mockClient.On("ListConsumableResources", mock.Anything, mock.Anything).
		Return(&batch.ListConsumableResourcesOutput{
			ConsumableResources: []batchtypes.ConsumableResourceSummary{
				{
					ConsumableResourceArn:  ptr.String("arn:aws:batch:us-east-1:123456789012:consumable-resource/my-cr"),
					ConsumableResourceName: ptr.String("my-cr"),
				},
			},
		}, nil)
	lister := &BatchConsumableResourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-cr", *resources[0].(*BatchConsumableResource).ConsumableResourceName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchConsumableResource_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)
	mockClient.On("ListConsumableResources", mock.Anything, mock.Anything).
		Return(&batch.ListConsumableResourcesOutput{ConsumableResources: []batchtypes.ConsumableResourceSummary{}}, nil)
	lister := &BatchConsumableResourceLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchConsumableResource_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)
	crArn := ptr.String("arn:aws:batch:us-east-1:123456789012:consumable-resource/my-cr")
	r := &BatchConsumableResource{
		svc:                   mockClient,
		ConsumableResourceArn: crArn,
	}
	mockClient.On("DeleteConsumableResource", mock.Anything,
		&batch.DeleteConsumableResourceInput{ConsumableResource: r.ConsumableResourceArn}).
		Return(&batch.DeleteConsumableResourceOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchConsumableResource_Properties(t *testing.T) {
	a := assert.New(t)
	r := BatchConsumableResource{
		ConsumableResourceArn:  ptr.String("arn"),
		ConsumableResourceName: ptr.String("my-cr"),
	}
	a.Equal("my-cr", r.Properties().Get("ConsumableResourceName"))
}

func Test_Mock_BatchConsumableResource_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-cr", (&BatchConsumableResource{ConsumableResourceName: ptr.String("my-cr")}).String())
}
