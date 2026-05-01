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

func Test_Mock_SageMakerMlflowApp_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListMlflowApps", mock.Anything, mock.Anything).
		Return(&sagemaker.ListMlflowAppsOutput{
			Summaries: []sagemakertypes.MlflowAppSummary{
				{
					Name:          ptr.String("my-mlflow-app"),
					Arn:           ptr.String("arn:aws:sagemaker:us-east-1:123456789012:mlflow-app/my-mlflow-app"),
					Status:        sagemakertypes.MlflowAppStatusCreated,
					MlflowVersion: ptr.String("2.0"),
				},
			},
		}, nil)

	lister := &SageMakerMlflowAppLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*SageMakerMlflowApp)
	a.Equal("my-mlflow-app", *r.Name)
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:mlflow-app/my-mlflow-app", *r.ARN)
	a.Equal(sagemakertypes.MlflowAppStatusCreated, r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerMlflowApp_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	mockClient.On("ListMlflowApps", mock.Anything, mock.Anything).
		Return(&sagemaker.ListMlflowAppsOutput{
			Summaries: []sagemakertypes.MlflowAppSummary{},
		}, nil)

	lister := &SageMakerMlflowAppLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerMlflowApp_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSageMakerV2Client)

	r := &SageMakerMlflowApp{
		svc:  mockClient,
		Name: ptr.String("my-mlflow-app"),
		ARN:  ptr.String("arn:aws:sagemaker:us-east-1:123456789012:mlflow-app/my-mlflow-app"),
	}

	mockClient.On("DeleteMlflowApp", mock.Anything,
		&sagemaker.DeleteMlflowAppInput{
			Arn: r.ARN,
		}).Return(&sagemaker.DeleteMlflowAppOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SageMakerMlflowApp_Properties(t *testing.T) {
	a := assert.New(t)
	r := &SageMakerMlflowApp{
		Name:          ptr.String("my-mlflow-app"),
		ARN:           ptr.String("arn:aws:sagemaker:us-east-1:123456789012:mlflow-app/my-mlflow-app"),
		Status:        sagemakertypes.MlflowAppStatusCreated,
		MlflowVersion: ptr.String("2.0"),
	}
	props := r.Properties()
	a.Equal("my-mlflow-app", props.Get("Name"))
	a.Equal("arn:aws:sagemaker:us-east-1:123456789012:mlflow-app/my-mlflow-app", props.Get("ARN"))
	a.Equal("Created", props.Get("Status"))
	a.Equal("2.0", props.Get("MlflowVersion"))
}

func Test_Mock_SageMakerMlflowApp_String(t *testing.T) {
	a := assert.New(t)
	r := &SageMakerMlflowApp{
		Name: ptr.String("my-mlflow-app"),
	}
	a.Equal("my-mlflow-app", r.String())
}

func Test_Mock_SageMakerMlflowApp_Filter_Deleting(t *testing.T) {
	a := assert.New(t)
	r := &SageMakerMlflowApp{
		Name:   ptr.String("my-mlflow-app"),
		Status: sagemakertypes.MlflowAppStatusDeleting,
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleting")
}

func Test_Mock_SageMakerMlflowApp_Filter_Deleted(t *testing.T) {
	a := assert.New(t)
	r := &SageMakerMlflowApp{
		Name:   ptr.String("my-mlflow-app"),
		Status: sagemakertypes.MlflowAppStatusDeleted,
	}
	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already deleted")
}

func Test_Mock_SageMakerMlflowApp_Filter_Active(t *testing.T) {
	a := assert.New(t)
	r := &SageMakerMlflowApp{
		Name:   ptr.String("my-mlflow-app"),
		Status: sagemakertypes.MlflowAppStatusCreated,
	}
	a.NoError(r.Filter())
}
