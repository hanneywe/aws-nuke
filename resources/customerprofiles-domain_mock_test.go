package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/customerprofiles"
	customerprofilestypes "github.com/aws/aws-sdk-go-v2/service/customerprofiles/types"
)

func Test_Mock_CustomerProfilesDomain_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListDomainsOutput{
			Items: []customerprofilestypes.ListDomainItem{
				{
					DomainName: ptr.String("my-domain"),
					Tags: map[string]string{
						"env": "test",
					},
				},
			},
		}, nil)

	lister := &CustomerProfilesDomainLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCustomerProfilesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	domain := resources[0].(*CustomerProfilesDomain)
	a.Equal("my-domain", *domain.DomainName)
	a.Equal("test", domain.Tags["env"])

	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesDomain_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	mockClient.
		On("ListDomains", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListDomainsOutput{
			Items: []customerprofilestypes.ListDomainItem{},
		}, nil)

	lister := &CustomerProfilesDomainLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCustomerProfilesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesDomain_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	domain := &CustomerProfilesDomain{
		svc:        mockClient,
		DomainName: ptr.String("my-domain"),
	}

	mockClient.
		On("DeleteDomain", mock.Anything, &customerprofiles.DeleteDomainInput{
			DomainName: domain.DomainName,
		}).
		Return(&customerprofiles.DeleteDomainOutput{}, nil)

	err := domain.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesDomain_Properties(t *testing.T) {
	a := assert.New(t)

	domain := CustomerProfilesDomain{
		DomainName: ptr.String("my-domain"),
		Tags: map[string]string{
			"env": "test",
		},
	}

	props := domain.Properties()
	a.Equal("my-domain", props.Get("DomainName"))
	a.Equal("test", props.Get("tag:env"))
}

func Test_Mock_CustomerProfilesDomain_String(t *testing.T) {
	a := assert.New(t)

	domain := CustomerProfilesDomain{
		DomainName: ptr.String("my-domain"),
	}

	a.Equal("my-domain", domain.String())
}
