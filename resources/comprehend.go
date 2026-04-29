package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/comprehend"
)

const ComprehendUnnamedJob = "Unnamed Job"

type ComprehendClient interface {
	ListFlywheels(ctx context.Context, params *comprehend.ListFlywheelsInput,
		optFns ...func(*comprehend.Options)) (*comprehend.ListFlywheelsOutput, error)
	DeleteFlywheel(ctx context.Context, params *comprehend.DeleteFlywheelInput,
		optFns ...func(*comprehend.Options)) (*comprehend.DeleteFlywheelOutput, error)
}
