package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codebuild"
)

type CodeBuildClient interface {
	ListFleets(ctx context.Context, params *codebuild.ListFleetsInput,
		optFns ...func(*codebuild.Options)) (*codebuild.ListFleetsOutput, error)
	BatchGetFleets(ctx context.Context, params *codebuild.BatchGetFleetsInput,
		optFns ...func(*codebuild.Options)) (*codebuild.BatchGetFleetsOutput, error)
	DeleteFleet(ctx context.Context, params *codebuild.DeleteFleetInput,
		optFns ...func(*codebuild.Options)) (*codebuild.DeleteFleetOutput, error)
}
