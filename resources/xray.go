package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/xray"
)

type XRayClient interface {
	GetEncryptionConfig(ctx context.Context, params *xray.GetEncryptionConfigInput,
		optFns ...func(*xray.Options)) (*xray.GetEncryptionConfigOutput, error)
	PutEncryptionConfig(ctx context.Context, params *xray.PutEncryptionConfigInput,
		optFns ...func(*xray.Options)) (*xray.PutEncryptionConfigOutput, error)
}
