package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	ecrtypes "github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

func Test_Mock_ECRImageScanningConfiguration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribeRepositories", mock.Anything, mock.Anything).
		Return(&ecr.DescribeRepositoriesOutput{
			Repositories: []ecrtypes.Repository{
				{
					RepositoryName:             ptr.String("repo-with-scan"),
					ImageScanningConfiguration: &ecrtypes.ImageScanningConfiguration{ScanOnPush: true},
				},
				{
					RepositoryName:             ptr.String("repo-without-scan"),
					ImageScanningConfiguration: &ecrtypes.ImageScanningConfiguration{ScanOnPush: false},
				},
			},
		}, nil)
	lister := &ECRImageScanningConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("repo-with-scan", resources[0].(*ECRImageScanningConfiguration).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRImageScanningConfiguration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribeRepositories", mock.Anything, mock.Anything).
		Return(&ecr.DescribeRepositoriesOutput{
			Repositories: []ecrtypes.Repository{},
		}, nil)
	lister := &ECRImageScanningConfigurationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRImageScanningConfiguration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	r := &ECRImageScanningConfiguration{
		svc:            mockClient,
		RepositoryName: ptr.String("repo-with-scan"),
		ScanOnPush:     ptr.Bool(true),
	}
	mockClient.On("PutImageScanningConfiguration", mock.Anything, &ecr.PutImageScanningConfigurationInput{
		RepositoryName: r.RepositoryName,
		ImageScanningConfiguration: &ecrtypes.ImageScanningConfiguration{
			ScanOnPush: false,
		},
	}).Return(&ecr.PutImageScanningConfigurationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRImageScanningConfiguration_Properties(t *testing.T) {
	a := assert.New(t)
	r := ECRImageScanningConfiguration{
		RepositoryName: ptr.String("my-repo"),
		ScanOnPush:     ptr.Bool(true),
	}
	props := r.Properties()
	a.Equal("my-repo", props.Get("RepositoryName"))
	a.Equal("true", props.Get("ScanOnPush"))
}

func Test_Mock_ECRImageScanningConfiguration_String(t *testing.T) {
	a := assert.New(t)
	r := &ECRImageScanningConfiguration{RepositoryName: ptr.String("my-repo")}
	a.Equal("my-repo", r.String())
}
