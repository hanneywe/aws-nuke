package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"
	mailmanagertypes "github.com/aws/aws-sdk-go-v2/service/mailmanager/types"
)

func Test_Mock_MailManagerRuleSet_List(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMailManagerClient)

	mockClient.
		On("ListRuleSets", mock.Anything, mock.Anything).
		Return(
			&mailmanager.ListRuleSetsOutput{
				RuleSets: []mailmanagertypes.RuleSet{
					{
						RuleSetId:   ptr.String("rs-12345"),
						RuleSetName: ptr.String("test-ruleset"),
					},
				},
			}, nil,
		)

	lister := &MailManagerRuleSetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	ruleSet := resources[0].(*MailManagerRuleSet)
	assertions.Equal("rs-12345", *ruleSet.RuleSetID)
	assertions.Equal("test-ruleset", *ruleSet.RuleSetName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerRuleSet_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMailManagerClient)

	mockClient.
		On("ListRuleSets", mock.Anything, mock.Anything).
		Return(
			&mailmanager.ListRuleSetsOutput{
				RuleSets: []mailmanagertypes.RuleSet{},
			}, nil,
		)

	lister := &MailManagerRuleSetLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testMailManagerListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerRuleSet_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockMailManagerClient)

	ruleSet := &MailManagerRuleSet{
		svc:         mockClient,
		RuleSetID:   ptr.String("rs-12345"),
		RuleSetName: ptr.String("test-ruleset"),
	}

	mockClient.
		On(
			"DeleteRuleSet",
			mock.Anything,
			&mailmanager.DeleteRuleSetInput{
				RuleSetId: ruleSet.RuleSetID,
			},
		).
		Return(&mailmanager.DeleteRuleSetOutput{}, nil)

	err := ruleSet.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_MailManagerRuleSet_Properties(t *testing.T) {
	assertions := assert.New(t)

	ruleSet := MailManagerRuleSet{
		RuleSetID:   ptr.String("rs-12345"),
		RuleSetName: ptr.String("test-ruleset"),
	}

	properties := ruleSet.Properties()

	assertions.Equal("rs-12345", properties.Get("RuleSetId"))
	assertions.Equal("test-ruleset", properties.Get("RuleSetName"))
}

func Test_Mock_MailManagerRuleSet_String(t *testing.T) {
	assertions := assert.New(t)

	ruleSet := MailManagerRuleSet{
		RuleSetID:   ptr.String("rs-12345"),
		RuleSetName: ptr.String("test-ruleset"),
	}

	assertions.Equal("test-ruleset", ruleSet.String())
}
