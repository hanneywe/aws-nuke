package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/applicationsignals"
)

type ApplicationSignalsClient interface {
	ListGroupingAttributeDefinitions(ctx context.Context, params *applicationsignals.ListGroupingAttributeDefinitionsInput,
		optFns ...func(*applicationsignals.Options)) (*applicationsignals.ListGroupingAttributeDefinitionsOutput, error)
	DeleteGroupingConfiguration(ctx context.Context, params *applicationsignals.DeleteGroupingConfigurationInput,
		optFns ...func(*applicationsignals.Options)) (*applicationsignals.DeleteGroupingConfigurationOutput, error)
}
