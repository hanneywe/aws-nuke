package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
)

type mockRoute53ResolverV2Client struct {
	mock.Mock
}

func (m *mockRoute53ResolverV2Client) ListFirewallRuleGroups(ctx context.Context, params *route53resolver.ListFirewallRuleGroupsInput,
	_ ...func(*route53resolver.Options)) (*route53resolver.ListFirewallRuleGroupsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53resolver.ListFirewallRuleGroupsOutput), args.Error(1)
}

func (m *mockRoute53ResolverV2Client) ListFirewallRules(ctx context.Context, params *route53resolver.ListFirewallRulesInput,
	_ ...func(*route53resolver.Options)) (*route53resolver.ListFirewallRulesOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53resolver.ListFirewallRulesOutput), args.Error(1)
}

func (m *mockRoute53ResolverV2Client) DeleteFirewallRule(ctx context.Context, params *route53resolver.DeleteFirewallRuleInput,
	_ ...func(*route53resolver.Options)) (*route53resolver.DeleteFirewallRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*route53resolver.DeleteFirewallRuleOutput), args.Error(1)
}
