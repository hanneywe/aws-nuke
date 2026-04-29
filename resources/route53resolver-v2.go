package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
)

// Route53ResolverV2Client is an interface for the AWS Route53 Resolver SDK v2 client methods
// used by Route53 Resolver sub-resources that need a separate v2 interface.
// This is separate from the existing Route53ResolverAPI interface which is used by other resources.
type Route53ResolverV2Client interface {
	ListFirewallRuleGroups(ctx context.Context, params *route53resolver.ListFirewallRuleGroupsInput,
		optFns ...func(*route53resolver.Options)) (*route53resolver.ListFirewallRuleGroupsOutput, error)
	ListFirewallRules(ctx context.Context, params *route53resolver.ListFirewallRulesInput,
		optFns ...func(*route53resolver.Options)) (*route53resolver.ListFirewallRulesOutput, error)
	DeleteFirewallRule(ctx context.Context, params *route53resolver.DeleteFirewallRuleInput,
		optFns ...func(*route53resolver.Options)) (*route53resolver.DeleteFirewallRuleOutput, error)
}
