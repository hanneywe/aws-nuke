package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/emr"
)

type EMRV2Client interface {
	GetBlockPublicAccessConfiguration(ctx context.Context, params *emr.GetBlockPublicAccessConfigurationInput,
		optFns ...func(*emr.Options)) (*emr.GetBlockPublicAccessConfigurationOutput, error)
	PutBlockPublicAccessConfiguration(ctx context.Context, params *emr.PutBlockPublicAccessConfigurationInput,
		optFns ...func(*emr.Options)) (*emr.PutBlockPublicAccessConfigurationOutput, error)
}
