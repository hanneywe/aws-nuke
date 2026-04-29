package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/connectcases"
	connectcasestypes "github.com/aws/aws-sdk-go-v2/service/connectcases/types"
)

func Test_Mock_ConnectCasesTemplate_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&connectcases.ListDomainsOutput{
			Domains: []connectcasestypes.DomainSummary{
				{
					DomainId: ptr.String("domain-12345"),
					Name:     ptr.String("my-domain"),
				},
			},
		}, nil)

	mockClient.
		On("ListTemplates", mock.Anything, mock.Anything).
		Return(&connectcases.ListTemplatesOutput{
			Templates: []connectcasestypes.TemplateSummary{
				{
					TemplateId: ptr.String("template-001"),
					Name:       ptr.String("my-template"),
				},
			},
		}, nil)

	lister := &ConnectCasesTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	tmpl := resources[0].(*ConnectCasesTemplate)
	a.Equal("domain-12345", *tmpl.DomainID)
	a.Equal("template-001", *tmpl.TemplateID)
	a.Equal("my-template", *tmpl.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesTemplate_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&connectcases.ListDomainsOutput{
			Domains: []connectcasestypes.DomainSummary{},
		}, nil)

	lister := &ConnectCasesTemplateLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesTemplate_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	tmpl := &ConnectCasesTemplate{
		svc:        mockClient,
		DomainID:   ptr.String("domain-12345"),
		TemplateID: ptr.String("template-001"),
		Name:       ptr.String("my-template"),
	}

	mockClient.
		On("DeleteTemplate", mock.Anything, &connectcases.DeleteTemplateInput{
			DomainId:   tmpl.DomainID,
			TemplateId: tmpl.TemplateID,
		}).
		Return(&connectcases.DeleteTemplateOutput{}, nil)

	err := tmpl.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesTemplate_Properties(t *testing.T) {
	a := assert.New(t)

	tmpl := ConnectCasesTemplate{
		DomainID:   ptr.String("domain-12345"),
		TemplateID: ptr.String("template-001"),
		Name:       ptr.String("my-template"),
	}

	props := tmpl.Properties()
	a.Equal("domain-12345", props.Get("DomainId"))
	a.Equal("template-001", props.Get("TemplateId"))
	a.Equal("my-template", props.Get("Name"))
}

func Test_Mock_ConnectCasesTemplate_String(t *testing.T) {
	a := assert.New(t)

	tmpl := ConnectCasesTemplate{
		TemplateID: ptr.String("template-001"),
	}

	a.Equal("template-001", tmpl.String())
}
