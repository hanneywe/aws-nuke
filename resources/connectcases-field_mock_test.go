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

func Test_Mock_ConnectCasesField_List_One(t *testing.T) {
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
		On("ListFields", mock.Anything, mock.Anything).
		Return(&connectcases.ListFieldsOutput{
			Fields: []connectcasestypes.FieldSummary{
				{
					FieldId:   ptr.String("field-001"),
					Name:      ptr.String("my-field"),
					Namespace: connectcasestypes.FieldNamespaceCustom,
				},
			},
		}, nil)

	lister := &ConnectCasesFieldLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	field := resources[0].(*ConnectCasesField)
	a.Equal("domain-12345", *field.DomainID)
	a.Equal("field-001", *field.FieldID)
	a.Equal("my-field", *field.Name)
	a.Equal(connectcasestypes.FieldNamespaceCustom, field.Namespace)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesField_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&connectcases.ListDomainsOutput{
			Domains: []connectcasestypes.DomainSummary{},
		}, nil)

	lister := &ConnectCasesFieldLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesField_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	field := &ConnectCasesField{
		svc:      mockClient,
		DomainID: ptr.String("domain-12345"),
		FieldID:  ptr.String("field-001"),
		Name:     ptr.String("my-field"),
	}

	mockClient.
		On("DeleteField", mock.Anything, &connectcases.DeleteFieldInput{
			DomainId: field.DomainID,
			FieldId:  field.FieldID,
		}).
		Return(&connectcases.DeleteFieldOutput{}, nil)

	err := field.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesField_Properties(t *testing.T) {
	a := assert.New(t)

	field := ConnectCasesField{
		DomainID:  ptr.String("domain-12345"),
		FieldID:   ptr.String("field-001"),
		Name:      ptr.String("my-field"),
		Namespace: connectcasestypes.FieldNamespaceCustom,
	}

	props := field.Properties()
	a.Equal("domain-12345", props.Get("DomainId"))
	a.Equal("field-001", props.Get("FieldId"))
	a.Equal("my-field", props.Get("Name"))
	a.Equal("Custom", props.Get("Namespace"))
}

func Test_Mock_ConnectCasesField_String(t *testing.T) {
	a := assert.New(t)

	field := ConnectCasesField{
		FieldID: ptr.String("field-001"),
	}

	a.Equal("field-001", field.String())
}
