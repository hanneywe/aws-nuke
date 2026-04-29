package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness"
)

// Route53RecoveryReadinessClient is the interface for the Route53 Recovery Readiness SDK client methods.
type Route53RecoveryReadinessClient interface {
	ListCells(ctx context.Context, params *route53recoveryreadiness.ListCellsInput,
		optFns ...func(*route53recoveryreadiness.Options)) (*route53recoveryreadiness.ListCellsOutput, error)
	DeleteCell(ctx context.Context, params *route53recoveryreadiness.DeleteCellInput,
		optFns ...func(*route53recoveryreadiness.Options)) (*route53recoveryreadiness.DeleteCellOutput, error)
	ListRecoveryGroups(ctx context.Context, params *route53recoveryreadiness.ListRecoveryGroupsInput,
		optFns ...func(*route53recoveryreadiness.Options)) (*route53recoveryreadiness.ListRecoveryGroupsOutput, error)
	DeleteRecoveryGroup(ctx context.Context, params *route53recoveryreadiness.DeleteRecoveryGroupInput,
		optFns ...func(*route53recoveryreadiness.Options)) (*route53recoveryreadiness.DeleteRecoveryGroupOutput, error)
}
