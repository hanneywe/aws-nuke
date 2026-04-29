package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/codedeploy"
)

type CodeDeployV2Client interface {
	ListOnPremisesInstances(ctx context.Context, params *codedeploy.ListOnPremisesInstancesInput,
		optFns ...func(*codedeploy.Options)) (*codedeploy.ListOnPremisesInstancesOutput, error)
	DeregisterOnPremisesInstance(ctx context.Context, params *codedeploy.DeregisterOnPremisesInstanceInput,
		optFns ...func(*codedeploy.Options)) (*codedeploy.DeregisterOnPremisesInstanceOutput, error)
}
