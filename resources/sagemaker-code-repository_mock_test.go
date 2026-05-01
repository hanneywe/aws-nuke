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

func Test_Mock_SageMakerCodeRepository_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListCodeRepositories", mock.Anything, mock.Anything).
		Return(&sagemaker.ListCodeRepositoriesOutput{
			CodeRepositorySummaryList: []sagemakertypes.CodeRepositorySummary{
				{
					CodeRepositoryName: ptr.String("my-repo"),
					CodeRepositoryArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:code-repository/my-repo"),
				},
			},
		}, nil)

	lister := &SageMakerCodeRepositoryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	repo := resources[0].(*SageMakerCodeRepository)
	a.Equal("my-repo", *repo.CodeRepositoryName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:code-repository/my-repo", *repo.CodeRepositoryArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerCodeRepository_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListCodeRepositories", mock.Anything, mock.Anything).
		Return(&sagemaker.ListCodeRepositoriesOutput{
			CodeRepositorySummaryList: []sagemakertypes.CodeRepositorySummary{},
		}, nil)

	lister := &SageMakerCodeRepositoryLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerCodeRepository_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	repo := &SageMakerCodeRepository{
		svc:                mockClient,
		CodeRepositoryName: ptr.String("my-repo"),
	}

	mockClient.On("DeleteCodeRepository", mock.Anything, &sagemaker.DeleteCodeRepositoryInput{
		CodeRepositoryName: repo.CodeRepositoryName,
	}).Return(&sagemaker.DeleteCodeRepositoryOutput{}, nil)

	a.NoError(repo.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerCodeRepository_Properties(t *testing.T) {
	a := assert.New(t)

	repo := SageMakerCodeRepository{
		CodeRepositoryName: ptr.String("my-repo"),
		CodeRepositoryArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:code-repository/my-repo"),
	}

	props := repo.Properties()
	a.Equal("my-repo", props.Get("CodeRepositoryName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:code-repository/my-repo", props.Get("CodeRepositoryArn"))
}

func Test_Mock_SageMakerCodeRepository_String(t *testing.T) {
	a := assert.New(t)
	repo := SageMakerCodeRepository{CodeRepositoryName: ptr.String("my-repo")}
	a.Equal("my-repo", repo.String())
}
