package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53resolver"
	route53resolvertypes "github.com/aws/aws-sdk-go-v2/service/route53resolver/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

var testRoute53ResolverV2ListerOpts = &nuke.ListerOpts{
	Logger: logrus.WithField("test", true),
}

func Test_Mock_Route53ResolverFirewallRule_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53ResolverV2Client)

	mockClient.On("ListFirewallRuleGroups", mock.Anything, mock.Anything).
		Return(&route53resolver.ListFirewallRuleGroupsOutput{
			FirewallRuleGroups: []route53resolvertypes.FirewallRuleGroupMetadata{
				{Id: ptr.String("rslvr-frg-123456")},
			},
		}, nil)

	mockClient.On("ListFirewallRules", mock.Anything, mock.Anything).
		Return(&route53resolver.ListFirewallRulesOutput{
			FirewallRules: []route53resolvertypes.FirewallRule{
				{
					FirewallRuleGroupId:  ptr.String("rslvr-frg-123456"),
					FirewallDomainListId: ptr.String("rslvr-fdl-abcdef"),
					Name:                 ptr.String("block-bad-domains"),
				},
				{
					FirewallRuleGroupId:  ptr.String("rslvr-frg-123456"),
					FirewallDomainListId: ptr.String("rslvr-fdl-ghijkl"),
					Name:                 ptr.String("allow-good-domains"),
				},
			},
		}, nil)

	lister := &Route53ResolverFirewallRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53ResolverV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	firewallRule := resources[0].(*Route53ResolverFirewallRuleV2)
	assertions.Equal("rslvr-frg-123456", *firewallRule.FirewallRuleGroupID)
	assertions.Equal("rslvr-fdl-abcdef", *firewallRule.FirewallDomainListID)
	assertions.Equal("block-bad-domains", *firewallRule.Name)

	secondRule := resources[1].(*Route53ResolverFirewallRuleV2)
	assertions.Equal("rslvr-fdl-ghijkl", *secondRule.FirewallDomainListID)
	assertions.Equal("allow-good-domains", *secondRule.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53ResolverFirewallRule_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53ResolverV2Client)

	mockClient.On("ListFirewallRuleGroups", mock.Anything, mock.Anything).
		Return(&route53resolver.ListFirewallRuleGroupsOutput{
			FirewallRuleGroups: []route53resolvertypes.FirewallRuleGroupMetadata{
				{Id: ptr.String("rslvr-frg-123456")},
			},
		}, nil)

	mockClient.On("ListFirewallRules", mock.Anything, mock.Anything).
		Return(&route53resolver.ListFirewallRulesOutput{
			FirewallRules: []route53resolvertypes.FirewallRule{},
		}, nil)

	lister := &Route53ResolverFirewallRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53ResolverV2ListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53ResolverFirewallRule_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRoute53ResolverV2Client)

	firewallRule := &Route53ResolverFirewallRuleV2{
		svc:                  mockClient,
		FirewallRuleGroupID:  ptr.String("rslvr-frg-123456"),
		FirewallDomainListID: ptr.String("rslvr-fdl-abcdef"),
		Name:                 ptr.String("block-bad-domains"),
	}

	mockClient.On("DeleteFirewallRule", mock.Anything, &route53resolver.DeleteFirewallRuleInput{
		FirewallRuleGroupId:  firewallRule.FirewallRuleGroupID,
		FirewallDomainListId: firewallRule.FirewallDomainListID,
	}).Return(&route53resolver.DeleteFirewallRuleOutput{}, nil)

	err := firewallRule.Remove(context.TODO())
	assertions.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53ResolverFirewallRule_Properties(t *testing.T) {
	assertions := assert.New(t)

	firewallRule := Route53ResolverFirewallRuleV2{
		FirewallRuleGroupID:  ptr.String("rslvr-frg-123456"),
		FirewallDomainListID: ptr.String("rslvr-fdl-abcdef"),
		Name:                 ptr.String("block-bad-domains"),
	}

	properties := firewallRule.Properties()
	assertions.Equal("rslvr-frg-123456", properties.Get("FirewallRuleGroupID"))
	assertions.Equal("rslvr-fdl-abcdef", properties.Get("FirewallDomainListID"))
	assertions.Equal("block-bad-domains", properties.Get("Name"))
}

func Test_Mock_Route53ResolverFirewallRule_String(t *testing.T) {
	assertions := assert.New(t)
	firewallRule := Route53ResolverFirewallRuleV2{Name: ptr.String("block-bad-domains")}
	assertions.Equal("block-bad-domains", firewallRule.String())
}
