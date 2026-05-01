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

func Test_Mock_SageMakerModelCard_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListModelCards", mock.Anything, mock.Anything).
		Return(&sagemaker.ListModelCardsOutput{
			ModelCardSummaries: []sagemakertypes.ModelCardSummary{
				{
					ModelCardName: ptr.String("my-model-card"),
					ModelCardArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:model-card/my-model-card"),
				},
			},
		}, nil)

	lister := &SageMakerModelCardLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	card := resources[0].(*SageMakerModelCard)
	a.Equal("my-model-card", *card.ModelCardName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:model-card/my-model-card", *card.ModelCardArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerModelCard_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListModelCards", mock.Anything, mock.Anything).
		Return(&sagemaker.ListModelCardsOutput{
			ModelCardSummaries: []sagemakertypes.ModelCardSummary{},
		}, nil)

	lister := &SageMakerModelCardLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerModelCard_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	card := &SageMakerModelCard{
		svc:           mockClient,
		ModelCardName: ptr.String("my-model-card"),
	}

	mockClient.On("DeleteModelCard", mock.Anything, &sagemaker.DeleteModelCardInput{
		ModelCardName: card.ModelCardName,
	}).Return(&sagemaker.DeleteModelCardOutput{}, nil)

	a.NoError(card.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerModelCard_Properties(t *testing.T) {
	a := assert.New(t)

	card := SageMakerModelCard{
		ModelCardName: ptr.String("my-model-card"),
		ModelCardArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:model-card/my-model-card"),
	}

	props := card.Properties()
	a.Equal("my-model-card", props.Get("ModelCardName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:model-card/my-model-card", props.Get("ModelCardArn"))
}

func Test_Mock_SageMakerModelCard_String(t *testing.T) {
	a := assert.New(t)
	card := SageMakerModelCard{ModelCardName: ptr.String("my-model-card")}
	a.Equal("my-model-card", card.String())
}
