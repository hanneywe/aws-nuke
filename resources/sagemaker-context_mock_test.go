package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	sagemakertypes "github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

func Test_Mock_SageMakerContext_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListContexts", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListContextsOutput{
				ContextSummaries: []sagemakertypes.ContextSummary{
					{
						ContextName: ptr.String("test-context"),
						ContextArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:context/test-context"),
					},
				},
			}, nil,
		)

	lister := &SageMakerContextLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	sagemakerContext := resources[0].(*SageMakerContext)
	assertions.Equal("test-context", *sagemakerContext.ContextName)
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:context/test-context", *sagemakerContext.ContextArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerContext_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.
		On("ListContexts", mock.Anything, mock.Anything).
		Return(
			&sagemaker.ListContextsOutput{
				ContextSummaries: []sagemakertypes.ContextSummary{},
			}, nil,
		)

	lister := &SageMakerContextLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerContext_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	sagemakerContext := &SageMakerContext{
		svc:         mockClient,
		ContextName: ptr.String("test-context"),
		ContextArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:context/test-context"),
	}

	mockClient.
		On("DeleteContext", mock.Anything,
			&sagemaker.DeleteContextInput{
				ContextName: sagemakerContext.ContextName,
			},
		).
		Return(&sagemaker.DeleteContextOutput{}, nil)

	err := sagemakerContext.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerContext_Properties(t *testing.T) {
	assertions := assert.New(t)

	sagemakerContext := SageMakerContext{
		ContextName: ptr.String("test-context"),
		ContextArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:context/test-context"),
	}

	properties := sagemakerContext.Properties()

	assertions.Equal("test-context", properties.Get("ContextName"))
	assertions.Equal("arn:aws:sagemaker:us-east-1:123456789012:context/test-context", properties.Get("ContextArn"))
}

func Test_Mock_SageMakerContext_String(t *testing.T) {
	assertions := assert.New(t)

	sagemakerContext := SageMakerContext{
		ContextName: ptr.String("test-context"),
	}

	assertions.Equal("test-context", sagemakerContext.String())
}
