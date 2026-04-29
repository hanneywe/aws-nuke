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

func Test_Mock_ConnectCasesDomain_List_One(t *testing.T) {
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

	lister := &ConnectCasesDomainLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	domain := resources[0].(*ConnectCasesDomain)
	a.Equal("domain-12345", *domain.DomainID)
	a.Equal("my-domain", *domain.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesDomain_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&connectcases.ListDomainsOutput{
			Domains: []connectcasestypes.DomainSummary{},
		}, nil)

	lister := &ConnectCasesDomainLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConnectCasesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesDomain_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConnectCasesClient)

	domain := &ConnectCasesDomain{
		svc:      mockClient,
		DomainID: ptr.String("domain-12345"),
		Name:     ptr.String("my-domain"),
	}

	mockClient.
		On("DeleteDomain", mock.Anything, &connectcases.DeleteDomainInput{
			DomainId: domain.DomainID,
		}).
		Return(&connectcases.DeleteDomainOutput{}, nil)

	err := domain.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_ConnectCasesDomain_Properties(t *testing.T) {
	a := assert.New(t)

	domain := ConnectCasesDomain{
		DomainID: ptr.String("domain-12345"),
		Name:     ptr.String("my-domain"),
	}

	props := domain.Properties()
	a.Equal("domain-12345", props.Get("DomainId"))
	a.Equal("my-domain", props.Get("Name"))
}

func Test_Mock_ConnectCasesDomain_String(t *testing.T) {
	a := assert.New(t)

	domain := ConnectCasesDomain{
		DomainID: ptr.String("domain-12345"),
	}

	a.Equal("domain-12345", domain.String())
}
