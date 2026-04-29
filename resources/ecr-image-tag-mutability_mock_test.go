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

func Test_Mock_ECRImageTagMutability_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribeRepositories", mock.Anything, mock.Anything).
		Return(&ecr.DescribeRepositoriesOutput{
			Repositories: []ecrtypes.Repository{
				{
					RepositoryName:     ptr.String("immutable-repo"),
					ImageTagMutability: ecrtypes.ImageTagMutabilityImmutable,
				},
				{
					RepositoryName:     ptr.String("mutable-repo"),
					ImageTagMutability: ecrtypes.ImageTagMutabilityMutable,
				},
			},
		}, nil)
	lister := &ECRImageTagMutabilityLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("immutable-repo", resources[0].(*ECRImageTagMutability).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRImageTagMutability_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribeRepositories", mock.Anything, mock.Anything).
		Return(&ecr.DescribeRepositoriesOutput{
			Repositories: []ecrtypes.Repository{},
		}, nil)
	lister := &ECRImageTagMutabilityLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRImageTagMutability_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	r := &ECRImageTagMutability{
		svc:                mockClient,
		RepositoryName:     ptr.String("immutable-repo"),
		ImageTagMutability: ptr.String("IMMUTABLE"),
	}
	mockClient.On("PutImageTagMutability", mock.Anything, &ecr.PutImageTagMutabilityInput{
		RepositoryName:     r.RepositoryName,
		ImageTagMutability: ecrtypes.ImageTagMutabilityMutable,
	}).Return(&ecr.PutImageTagMutabilityOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRImageTagMutability_Properties(t *testing.T) {
	a := assert.New(t)
	r := ECRImageTagMutability{
		RepositoryName:     ptr.String("my-repo"),
		ImageTagMutability: ptr.String("IMMUTABLE"),
	}
	props := r.Properties()
	a.Equal("my-repo", props.Get("RepositoryName"))
	a.Equal("IMMUTABLE", props.Get("ImageTagMutability"))
}

func Test_Mock_ECRImageTagMutability_String(t *testing.T) {
	a := assert.New(t)
	r := &ECRImageTagMutability{RepositoryName: ptr.String("my-repo")}
	a.Equal("my-repo", r.String())
}
