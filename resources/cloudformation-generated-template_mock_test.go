package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

func Test_Mock_CloudFormationGeneratedTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudFormationClient)
	mockClient.On("ListGeneratedTemplates", mock.Anything, mock.Anything).
		Return(&cloudformation.ListGeneratedTemplatesOutput{
			Summaries: []cfntypes.TemplateSummary{
				{GeneratedTemplateId: ptr.String("gt-1"), GeneratedTemplateName: ptr.String("my-template")},
			},
		}, nil)
	lister := &CloudFormationGeneratedTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudFormationListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-template", *resources[0].(*CloudFormationGeneratedTemplate).GeneratedTemplateName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudFormationGeneratedTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudFormationClient)
	mockClient.On("ListGeneratedTemplates", mock.Anything, mock.Anything).
		Return(&cloudformation.ListGeneratedTemplatesOutput{Summaries: []cfntypes.TemplateSummary{}}, nil)
	lister := &CloudFormationGeneratedTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCloudFormationListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudFormationGeneratedTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudFormationClient)
	r := &CloudFormationGeneratedTemplate{
		svc:                   mockClient,
		GeneratedTemplateName: ptr.String("my-template"),
	}
	mockClient.On("DeleteGeneratedTemplate", mock.Anything,
		&cloudformation.DeleteGeneratedTemplateInput{
			GeneratedTemplateName: r.GeneratedTemplateName,
		}).
		Return(&cloudformation.DeleteGeneratedTemplateOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudFormationGeneratedTemplate_Properties(t *testing.T) {
	a := assert.New(t)
	r := CloudFormationGeneratedTemplate{GeneratedTemplateID: ptr.String("gt-1"), GeneratedTemplateName: ptr.String("my-template")}
	a.Equal("my-template", r.Properties().Get("GeneratedTemplateName"))
	a.Equal("gt-1", r.Properties().Get("GeneratedTemplateId"))
}

func Test_Mock_CloudFormationGeneratedTemplate_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-template", (&CloudFormationGeneratedTemplate{GeneratedTemplateName: ptr.String("my-template")}).String())
}
