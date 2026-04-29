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

func Test_Mock_SESv2EmailIdentity_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListEmailIdentities", mock.Anything, mock.Anything).
		Return(&sesv2.ListEmailIdentitiesOutput{
			EmailIdentities: []sesv2types.IdentityInfo{
				{IdentityName: ptr.String("example.com"), IdentityType: sesv2types.IdentityTypeDomain},
			},
		}, nil)
	lister := &SESv2EmailIdentityLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	id := resources[0].(*SESv2EmailIdentity)
	a.Equal("example.com", *id.IdentityName)
	a.Equal("DOMAIN", *id.IdentityType)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2EmailIdentity_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	mockClient.On("ListEmailIdentities", mock.Anything, mock.Anything).
		Return(&sesv2.ListEmailIdentitiesOutput{EmailIdentities: []sesv2types.IdentityInfo{}}, nil)
	lister := &SESv2EmailIdentityLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSESv2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2EmailIdentity_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESv2Client)
	id := &SESv2EmailIdentity{svc: mockClient, IdentityName: ptr.String("example.com")}
	mockClient.On("DeleteEmailIdentity", mock.Anything, &sesv2.DeleteEmailIdentityInput{EmailIdentity: id.IdentityName}).
		Return(&sesv2.DeleteEmailIdentityOutput{}, nil)
	a.NoError(id.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESv2EmailIdentity_Properties(t *testing.T) {
	a := assert.New(t)
	id := SESv2EmailIdentity{IdentityName: ptr.String("example.com"), IdentityType: ptr.String("DOMAIN")}
	a.Equal("example.com", id.Properties().Get("IdentityName"))
	a.Equal("DOMAIN", id.Properties().Get("IdentityType"))
}

func Test_Mock_SESv2EmailIdentity_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("example.com", (&SESv2EmailIdentity{IdentityName: ptr.String("example.com")}).String())
}
