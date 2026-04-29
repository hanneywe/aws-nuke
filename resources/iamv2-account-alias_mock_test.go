package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func Test_Mock_IAMAccountAlias_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIAMClient)

	mockClient.On("ListAccountAliases", mock.Anything, mock.Anything).
		Return(&iam.ListAccountAliasesOutput{
			AccountAliases: []string{
				"test-alias",
			},
		}, nil)

	lister := &IAMAccountAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIamListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*IAMAccountAlias)
	a.Equal("test-alias", *r.AccountAlias)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IAMAccountAlias_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIAMClient)

	mockClient.On("ListAccountAliases", mock.Anything, mock.Anything).
		Return(&iam.ListAccountAliasesOutput{
			AccountAliases: []string{},
		}, nil)

	lister := &IAMAccountAliasLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIamListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_IAMAccountAlias_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockIAMClient)

	r := &IAMAccountAlias{
		svc:          mockClient,
		AccountAlias: ptr.String("test-accountalias"),
	}

	mockClient.On("DeleteAccountAlias", mock.Anything,
		&iam.DeleteAccountAliasInput{
			AccountAlias: r.AccountAlias,
		}).Return(&iam.DeleteAccountAliasOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IAMAccountAlias_Properties(t *testing.T) {
	a := assert.New(t)
	r := &IAMAccountAlias{
		AccountAlias: ptr.String("test-accountalias"),
	}
	props := r.Properties()
	a.Equal("test-accountalias", props.Get("AccountAlias"))
}

func Test_Mock_IAMAccountAlias_String(t *testing.T) {
	a := assert.New(t)
	r := &IAMAccountAlias{
		AccountAlias: ptr.String("test-accountalias"),
	}
	a.Equal("test-accountalias", r.String())
}
