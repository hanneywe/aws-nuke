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

func Test_Mock_SageMakerModelPackageGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListModelPackageGroups", mock.Anything, mock.Anything).
		Return(&sagemaker.ListModelPackageGroupsOutput{
			ModelPackageGroupSummaryList: []sagemakertypes.ModelPackageGroupSummary{
				{
					ModelPackageGroupName: ptr.String("my-model-pkg-group"),
					ModelPackageGroupArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:model-package-group/my-model-pkg-group"),
				},
			},
		}, nil)

	lister := &SageMakerModelPackageGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	group := resources[0].(*SageMakerModelPackageGroup)
	a.Equal("my-model-pkg-group", *group.ModelPackageGroupName)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:model-package-group/my-model-pkg-group", *group.ModelPackageGroupArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerModelPackageGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListModelPackageGroups", mock.Anything, mock.Anything).
		Return(&sagemaker.ListModelPackageGroupsOutput{
			ModelPackageGroupSummaryList: []sagemakertypes.ModelPackageGroupSummary{},
		}, nil)

	lister := &SageMakerModelPackageGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSageMakerV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerModelPackageGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	group := &SageMakerModelPackageGroup{
		svc:                   mockClient,
		ModelPackageGroupName: ptr.String("my-model-pkg-group"),
	}

	mockClient.On("DeleteModelPackageGroup", mock.Anything, &sagemaker.DeleteModelPackageGroupInput{
		ModelPackageGroupName: group.ModelPackageGroupName,
	}).Return(&sagemaker.DeleteModelPackageGroupOutput{}, nil)

	a.NoError(group.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerModelPackageGroup_Properties(t *testing.T) {
	a := assert.New(t)

	group := SageMakerModelPackageGroup{
		ModelPackageGroupName: ptr.String("my-model-pkg-group"),
		ModelPackageGroupArn:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:model-package-group/my-model-pkg-group"),
	}

	props := group.Properties()
	a.Equal("my-model-pkg-group", props.Get("ModelPackageGroupName"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:model-package-group/my-model-pkg-group", props.Get("ModelPackageGroupArn"))
}

func Test_Mock_SageMakerModelPackageGroup_String(t *testing.T) {
	a := assert.New(t)
	group := SageMakerModelPackageGroup{ModelPackageGroupName: ptr.String("my-model-pkg-group")}
	a.Equal("my-model-pkg-group", group.String())
}
