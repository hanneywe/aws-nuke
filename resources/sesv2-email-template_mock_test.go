package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	sesv2types "github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

func Test_Mock_SESv2EmailTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListEmailTemplates", mock.Anything, mock.Anything).
		Return(&sesv2.ListEmailTemplatesOutput{
			TemplatesMetadata: []sesv2types.EmailTemplateMetadata{
				{TemplateName: ptr.String("my-template")},
			},
		}, nil)
	lister := &SESv2EmailTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	tmpl := resources[0].(*SESv2EmailTemplate)
	a.Equal("my-template", *tmpl.TemplateName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2EmailTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListEmailTemplates", mock.Anything, mock.Anything).
		Return(&sesv2.ListEmailTemplatesOutput{TemplatesMetadata: []sesv2types.EmailTemplateMetadata{}}, nil)
	lister := &SESv2EmailTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2EmailTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	tmpl := &SESv2EmailTemplate{svc: mockClient, TemplateName: ptr.String("my-template")}
	mockClient.On("DeleteEmailTemplate", mock.Anything, &sesv2.DeleteEmailTemplateInput{TemplateName: tmpl.TemplateName}).
		Return(&sesv2.DeleteEmailTemplateOutput{}, nil)
	a.NoError(tmpl.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2EmailTemplate_Properties(t *testing.T) {
	a := assert.New(t)
	tmpl := SESv2EmailTemplate{TemplateName: ptr.String("my-template")}
	a.Equal("my-template", tmpl.Properties().Get("TemplateName"))
}

func Test_Mock_SESv2EmailTemplate_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-template", (&SESv2EmailTemplate{TemplateName: ptr.String("my-template")}).String())
}
