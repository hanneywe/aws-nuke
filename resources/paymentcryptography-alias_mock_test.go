package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/paymentcryptography"
	paymentcryptographytypes "github.com/aws/aws-sdk-go-v2/service/paymentcryptography/types"
)

func Test_Mock_PaymentCryptographyAlias_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockPaymentCryptographyClient)

	mockClient.
		On("ListAliases", mock.Anything, mock.Anything).
		Return(&paymentcryptography.ListAliasesOutput{
			Aliases: []paymentcryptographytypes.Alias{
				{
					AliasName: ptr.String("alias/test-alias"),
					KeyArn:    ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
				},
			},
		}, nil)

	lister := &PaymentCryptographyAliasLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testPaymentCryptographyListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	alias := resources[0].(*PaymentCryptographyAlias)
	assertions.Equal("alias/test-alias", *alias.AliasName)
	assertions.Equal("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123", *alias.KeyArn)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PaymentCryptographyAlias_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockPaymentCryptographyClient)

	mockClient.
		On("ListAliases", mock.Anything, mock.Anything).
		Return(&paymentcryptography.ListAliasesOutput{
			Aliases: []paymentcryptographytypes.Alias{},
		}, nil)

	lister := &PaymentCryptographyAliasLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testPaymentCryptographyListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PaymentCryptographyAlias_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockPaymentCryptographyClient)

	alias := &PaymentCryptographyAlias{
		svc:       mockClient,
		AliasName: ptr.String("alias/test-alias"),
	}

	mockClient.
		On("DeleteAlias", mock.Anything, &paymentcryptography.DeleteAliasInput{
			AliasName: alias.AliasName,
		}).
		Return(&paymentcryptography.DeleteAliasOutput{}, nil)

	err := alias.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PaymentCryptographyAlias_Properties(t *testing.T) {
	assertions := assert.New(t)

	alias := PaymentCryptographyAlias{
		AliasName: ptr.String("alias/test-alias"),
		KeyArn:    ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
	}

	properties := alias.Properties()

	assertions.Equal("alias/test-alias", properties.Get("AliasName"))
	assertions.Equal("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123", properties.Get("KeyArn"))
}

func Test_Mock_PaymentCryptographyAlias_String(t *testing.T) {
	assertions := assert.New(t)

	alias := PaymentCryptographyAlias{
		AliasName: ptr.String("alias/test-alias"),
	}

	assertions.Equal("alias/test-alias", alias.String())
}
