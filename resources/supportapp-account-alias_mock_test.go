package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/supportapp"
)

func Test_Mock_SupportAppAccountAlias_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSupportAppClient)
	mockClient.On("GetAccountAlias", mock.Anything, mock.Anything).
		Return(&supportapp.GetAccountAliasOutput{
			AccountAlias: ptr.String("my-account-alias"),
		}, nil)
	lister := &SupportAppAccountAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSupportAppListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	r := resources[0].(*SupportAppAccountAlias)
	a.Equal("my-account-alias", *r.AccountAlias)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SupportAppAccountAlias_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSupportAppClient)
	mockClient.On("GetAccountAlias", mock.Anything, mock.Anything).
		Return(&supportapp.GetAccountAliasOutput{
			AccountAlias: ptr.String(""),
		}, nil)
	lister := &SupportAppAccountAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSupportAppListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SupportAppAccountAlias_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSupportAppClient)
	r := &SupportAppAccountAlias{svc: mockClient, AccountAlias: ptr.String("my-account-alias")}
	mockClient.On("DeleteAccountAlias", mock.Anything, &supportapp.DeleteAccountAliasInput{}).
		Return(&supportapp.DeleteAccountAliasOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SupportAppAccountAlias_Properties(t *testing.T) {
	a := assert.New(t)
	r := SupportAppAccountAlias{
		AccountAlias: ptr.String("my-account-alias"),
	}
	props := r.Properties()
	a.Equal("my-account-alias", props.Get("AccountAlias"))
}

func Test_Mock_SupportAppAccountAlias_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-account-alias", (&SupportAppAccountAlias{AccountAlias: ptr.String("my-account-alias")}).String())
}
