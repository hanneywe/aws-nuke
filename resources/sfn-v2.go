package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type SFNv2Client interface {
	ListActivities(ctx context.Context, params *sfn.ListActivitiesInput,
		optFns ...func(*sfn.Options)) (*sfn.ListActivitiesOutput, error)
	DeleteActivity(ctx context.Context, params *sfn.DeleteActivityInput,
		optFns ...func(*sfn.Options)) (*sfn.DeleteActivityOutput, error)
}
