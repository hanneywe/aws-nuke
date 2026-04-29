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

func Test_Mock_ECRRepositoryCreationTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribeRepositoryCreationTemplates", mock.Anything, mock.Anything).
		Return(&ecr.DescribeRepositoryCreationTemplatesOutput{
			RepositoryCreationTemplates: []ecrtypes.RepositoryCreationTemplate{
				{Prefix: ptr.String("my-prefix")},
			},
		}, nil)
	lister := &ECRRepositoryCreationTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-prefix", resources[0].(*ECRRepositoryCreationTemplate).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRRepositoryCreationTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	mockClient.On("DescribeRepositoryCreationTemplates", mock.Anything, mock.Anything).
		Return(&ecr.DescribeRepositoryCreationTemplatesOutput{RepositoryCreationTemplates: []ecrtypes.RepositoryCreationTemplate{}}, nil)
	lister := &ECRRepositoryCreationTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testECRv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRRepositoryCreationTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockECRv2Client)
	r := &ECRRepositoryCreationTemplate{svc: mockClient, Prefix: ptr.String("my-prefix")}
	mockClient.On("DeleteRepositoryCreationTemplate", mock.Anything, &ecr.DeleteRepositoryCreationTemplateInput{
		Prefix: r.Prefix,
	}).Return(&ecr.DeleteRepositoryCreationTemplateOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ECRRepositoryCreationTemplate_Properties(t *testing.T) {
	a := assert.New(t)
	r := ECRRepositoryCreationTemplate{Prefix: ptr.String("my-prefix")}
	a.Equal("my-prefix", r.Properties().Get("Prefix"))
}

func Test_Mock_ECRRepositoryCreationTemplate_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-prefix", (&ECRRepositoryCreationTemplate{Prefix: ptr.String("my-prefix")}).String())
}
