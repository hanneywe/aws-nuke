package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/paymentcryptography"
	paymentcryptographytypes "github.com/aws/aws-sdk-go-v2/service/paymentcryptography/types"
)

func Test_Mock_PaymentCryptographyKey_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockPaymentCryptographyClient)

	mockClient.
		On("ListKeys", mock.Anything, mock.Anything).
		Return(&paymentcryptography.ListKeysOutput{
			Keys: []paymentcryptographytypes.KeySummary{
				{
					KeyArn:   ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
					KeyState: paymentcryptographytypes.KeyStateCreateComplete,
				},
			},
		}, nil)

	lister := &PaymentCryptographyKeyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testPaymentCryptographyListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	paymentKey := resources[0].(*PaymentCryptographyKey)
	assertions.Equal("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123", *paymentKey.KeyArn)
	assertions.Equal(paymentcryptographytypes.KeyStateCreateComplete, paymentKey.KeyState)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PaymentCryptographyKey_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockPaymentCryptographyClient)

	mockClient.
		On("ListKeys", mock.Anything, mock.Anything).
		Return(&paymentcryptography.ListKeysOutput{
			Keys: []paymentcryptographytypes.KeySummary{},
		}, nil)

	lister := &PaymentCryptographyKeyLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testPaymentCryptographyListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PaymentCryptographyKey_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockPaymentCryptographyClient)

	paymentKey := &PaymentCryptographyKey{
		svc:    mockClient,
		KeyArn: ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
	}

	mockClient.
		On("DeleteKey", mock.Anything, &paymentcryptography.DeleteKeyInput{
			KeyIdentifier:   paymentKey.KeyArn,
			DeleteKeyInDays: aws.Int32(3),
		}).
		Return(&paymentcryptography.DeleteKeyOutput{}, nil)

	err := paymentKey.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_PaymentCryptographyKey_Properties(t *testing.T) {
	assertions := assert.New(t)

	paymentKey := PaymentCryptographyKey{
		KeyArn: ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
	}

	properties := paymentKey.Properties()

	assertions.Equal("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123", properties.Get("KeyArn"))
}

func Test_Mock_PaymentCryptographyKey_String(t *testing.T) {
	assertions := assert.New(t)

	paymentKey := PaymentCryptographyKey{
		KeyArn: ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
	}

	assertions.Equal("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123", paymentKey.String())
}

func Test_Mock_PaymentCryptographyKey_Filter_Active(t *testing.T) {
	assertions := assert.New(t)

	paymentKey := PaymentCryptographyKey{
		KeyArn:   ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
		KeyState: paymentcryptographytypes.KeyStateCreateComplete,
	}

	err := paymentKey.Filter()
	assertions.NoError(err)
}

func Test_Mock_PaymentCryptographyKey_Filter_DeletePending(t *testing.T) {
	assertions := assert.New(t)

	paymentKey := PaymentCryptographyKey{
		KeyArn:   ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
		KeyState: paymentcryptographytypes.KeyStateDeletePending,
	}

	err := paymentKey.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETE_PENDING")
}

func Test_Mock_PaymentCryptographyKey_Filter_DeleteComplete(t *testing.T) {
	assertions := assert.New(t)

	paymentKey := PaymentCryptographyKey{
		KeyArn:   ptr.String("arn:aws:payment-cryptography:us-east-1:123456789012:key/abc123"),
		KeyState: paymentcryptographytypes.KeyStateDeleteComplete,
	}

	err := paymentKey.Filter()
	assertions.Error(err)
	assertions.Contains(err.Error(), "DELETE_COMPLETE")
}
