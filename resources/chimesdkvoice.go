package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkvoice"
)

type ChimeSDKVoiceClient interface {
	ListPhoneNumbers(ctx context.Context, params *chimesdkvoice.ListPhoneNumbersInput,
		optFns ...func(*chimesdkvoice.Options)) (*chimesdkvoice.ListPhoneNumbersOutput, error)
	DeletePhoneNumber(ctx context.Context, params *chimesdkvoice.DeletePhoneNumberInput,
		optFns ...func(*chimesdkvoice.Options)) (*chimesdkvoice.DeletePhoneNumberOutput, error)
}
