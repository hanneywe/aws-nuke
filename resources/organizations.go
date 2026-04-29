package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
)

type OrganizationsClient interface {
	ListAccounts(ctx context.Context, params *organizations.ListAccountsInput,
		optFns ...func(*organizations.Options)) (*organizations.ListAccountsOutput, error)
	CloseAccount(ctx context.Context, params *organizations.CloseAccountInput,
		optFns ...func(*organizations.Options)) (*organizations.CloseAccountOutput, error)
	DescribeOrganization(ctx context.Context, params *organizations.DescribeOrganizationInput,
		optFns ...func(*organizations.Options)) (*organizations.DescribeOrganizationOutput, error)
}
