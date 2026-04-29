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

func Test_Mock_CustomerProfilesProfile_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListDomainsOutput{
			Items: []customerprofilestypes.ListDomainItem{
				{DomainName: ptr.String("my-domain")},
			},
		}, nil)

	mockClient.On("SearchProfiles", mock.Anything, mock.Anything).
		Return(&customerprofiles.SearchProfilesOutput{
			Items: []customerprofilestypes.Profile{
				{
					ProfileId: ptr.String("profile-123"),
					FirstName: ptr.String("John"),
					LastName:  ptr.String("Doe"),
				},
			},
		}, nil)

	lister := &CustomerProfilesProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCustomerProfilesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CustomerProfilesProfile)
	a.Equal("my-domain", *r.DomainName)
	a.Equal("profile-123", *r.ProfileID)
	a.Equal("John", *r.FirstName)
	a.Equal("Doe", *r.LastName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesProfile_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	mockClient.On("ListDomains", mock.Anything, mock.Anything).
		Return(&customerprofiles.ListDomainsOutput{
			Items: []customerprofilestypes.ListDomainItem{},
		}, nil)

	lister := &CustomerProfilesProfileLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testCustomerProfilesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesProfile_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCustomerProfilesClient)

	r := &CustomerProfilesProfile{
		svc:        mockClient,
		DomainName: ptr.String("my-domain"),
		ProfileID:  ptr.String("profile-123"),
	}

	mockClient.On("DeleteProfile", mock.Anything,
		&customerprofiles.DeleteProfileInput{
			DomainName: r.DomainName,
			ProfileId:  r.ProfileID,
		}).Return(&customerprofiles.DeleteProfileOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_CustomerProfilesProfile_Properties(t *testing.T) {
	a := assert.New(t)
	r := &CustomerProfilesProfile{
		DomainName: ptr.String("my-domain"),
		ProfileID:  ptr.String("profile-123"),
		FirstName:  ptr.String("John"),
		LastName:   ptr.String("Doe"),
	}
	props := r.Properties()
	a.Equal("my-domain", props.Get("DomainName"))
	a.Equal("profile-123", props.Get("ProfileID"))
	a.Equal("John", props.Get("FirstName"))
	a.Equal("Doe", props.Get("LastName"))
}

func Test_Mock_CustomerProfilesProfile_String(t *testing.T) {
	a := assert.New(t)
	r := &CustomerProfilesProfile{
		ProfileID: ptr.String("profile-123"),
	}
	a.Equal("profile-123", r.String())
}
