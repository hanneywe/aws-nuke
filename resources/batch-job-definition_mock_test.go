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

func Test_Mock_BatchJobDefinition_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)

	mockClient.On("DescribeJobDefinitions", mock.Anything, mock.Anything).
		Return(&batch.DescribeJobDefinitionsOutput{
			JobDefinitions: []batchtypes.JobDefinition{
				{
					JobDefinitionArn:  ptr.String("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1"),
					JobDefinitionName: ptr.String("my-job"),
					Status:            ptr.String("ACTIVE"),
				},
			},
		}, nil)

	lister := &BatchJobDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*BatchJobDefinition)
	a.Equal("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1", *r.JobDefinitionArn)
	a.Equal("my-job", *r.JobDefinitionName)
	a.Equal("ACTIVE", *r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchJobDefinition_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)

	mockClient.On("DescribeJobDefinitions", mock.Anything, mock.Anything).
		Return(&batch.DescribeJobDefinitionsOutput{
			JobDefinitions: []batchtypes.JobDefinition{},
		}, nil)

	lister := &BatchJobDefinitionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testBatchListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchJobDefinition_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockBatchClient)

	r := &BatchJobDefinition{
		svc:              mockClient,
		JobDefinitionArn: ptr.String("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1"),
	}

	mockClient.On("DeregisterJobDefinition", mock.Anything,
		&batch.DeregisterJobDefinitionInput{
			JobDefinition: r.JobDefinitionArn,
		}).Return(&batch.DeregisterJobDefinitionOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_BatchJobDefinition_Properties(t *testing.T) {
	a := assert.New(t)
	r := &BatchJobDefinition{
		JobDefinitionArn:  ptr.String("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1"),
		JobDefinitionName: ptr.String("my-job"),
		Status:            ptr.String("ACTIVE"),
	}
	props := r.Properties()
	a.Equal("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1", props.Get("JobDefinitionArn"))
	a.Equal("my-job", props.Get("JobDefinitionName"))
	a.Equal("ACTIVE", props.Get("Status"))
}

func Test_Mock_BatchJobDefinition_String(t *testing.T) {
	a := assert.New(t)
	r := &BatchJobDefinition{
		JobDefinitionArn: ptr.String("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1"),
	}
	a.Equal("arn:aws:batch:us-east-1:123456789012:job-definition/my-job:1", r.String())
}

func Test_Mock_BatchJobDefinition_Filter_Inactive(t *testing.T) {
	a := assert.New(t)
	r := &BatchJobDefinition{
		Status: ptr.String("INACTIVE"),
	}
	a.Error(r.Filter())
}

func Test_Mock_BatchJobDefinition_Filter_Active(t *testing.T) {
	a := assert.New(t)
	r := &BatchJobDefinition{
		Status: ptr.String("ACTIVE"),
	}
	a.NoError(r.Filter())
}
