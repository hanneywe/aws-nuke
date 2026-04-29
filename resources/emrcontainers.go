package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/emrcontainers"
)

// EMRContainersClient is the interface for the EMR Containers SDK client methods.
type EMRContainersClient interface {
	ListJobTemplates(ctx context.Context, params *emrcontainers.ListJobTemplatesInput,
		optFns ...func(*emrcontainers.Options)) (*emrcontainers.ListJobTemplatesOutput, error)
	DeleteJobTemplate(ctx context.Context, params *emrcontainers.DeleteJobTemplateInput,
		optFns ...func(*emrcontainers.Options)) (*emrcontainers.DeleteJobTemplateOutput, error)
}
