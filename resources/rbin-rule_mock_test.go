package resources

import (
	"context"
	"testing"
	"time"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/rbin"
	rbintypes "github.com/aws/aws-sdk-go-v2/service/rbin/types"

	libsettings "github.com/ekristen/libnuke/pkg/settings"
)

func Test_Mock_RbinRule_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	// First call for EBS_SNAPSHOT returns one rule
	mockClient.
		On("ListRules", mock.Anything, mock.MatchedBy(func(input *rbin.ListRulesInput) bool {
			return input.ResourceType == rbintypes.ResourceTypeEbsSnapshot
		})).
		Return(&rbin.ListRulesOutput{
			Rules: []rbintypes.RuleSummary{
				{
					Identifier:  ptr.String("rule-abc123"),
					Description: ptr.String("Retain EBS snapshots"),
					LockState:   rbintypes.LockStateLocked,
				},
			},
		}, nil)

	// Second call for EC2_IMAGE returns empty
	mockClient.
		On("ListRules", mock.Anything, mock.MatchedBy(func(input *rbin.ListRulesInput) bool {
			return input.ResourceType == rbintypes.ResourceTypeEc2Image
		})).
		Return(&rbin.ListRulesOutput{
			Rules: []rbintypes.RuleSummary{},
		}, nil)

	lister := &RbinRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testRbinListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	rbinRule := resources[0].(*RbinRule)
	assertions.Equal("rule-abc123", *rbinRule.Identifier)
	assertions.Equal("Retain EBS snapshots", *rbinRule.Description)
	assertions.Equal("EBS_SNAPSHOT", rbinRule.ResourceType)
	assertions.Equal("locked", rbinRule.LockState)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	// Both resource types return empty
	mockClient.
		On("ListRules", mock.Anything, mock.MatchedBy(func(input *rbin.ListRulesInput) bool {
			return input.ResourceType == rbintypes.ResourceTypeEbsSnapshot
		})).
		Return(&rbin.ListRulesOutput{
			Rules: []rbintypes.RuleSummary{},
		}, nil)

	mockClient.
		On("ListRules", mock.Anything, mock.MatchedBy(func(input *rbin.ListRulesInput) bool {
			return input.ResourceType == rbintypes.ResourceTypeEc2Image
		})).
		Return(&rbin.ListRulesOutput{
			Rules: []rbintypes.RuleSummary{},
		}, nil)

	lister := &RbinRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testRbinListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_List_MultipleResourceTypes(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	// EBS_SNAPSHOT returns one rule
	mockClient.
		On("ListRules", mock.Anything, mock.MatchedBy(func(input *rbin.ListRulesInput) bool {
			return input.ResourceType == rbintypes.ResourceTypeEbsSnapshot
		})).
		Return(&rbin.ListRulesOutput{
			Rules: []rbintypes.RuleSummary{
				{
					Identifier:  ptr.String("rule-snap-1"),
					Description: ptr.String("Snapshot rule"),
				},
			},
		}, nil)

	// EC2_IMAGE returns one rule
	mockClient.
		On("ListRules", mock.Anything, mock.MatchedBy(func(input *rbin.ListRulesInput) bool {
			return input.ResourceType == rbintypes.ResourceTypeEc2Image
		})).
		Return(&rbin.ListRulesOutput{
			Rules: []rbintypes.RuleSummary{
				{
					Identifier:  ptr.String("rule-ami-1"),
					Description: ptr.String("AMI rule"),
				},
			},
		}, nil)

	lister := &RbinRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testRbinListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 2)

	snapshotRule := resources[0].(*RbinRule)
	assertions.Equal("rule-snap-1", *snapshotRule.Identifier)
	assertions.Equal("EBS_SNAPSHOT", snapshotRule.ResourceType)

	imageRule := resources[1].(*RbinRule)
	assertions.Equal("rule-ami-1", *imageRule.Identifier)
	assertions.Equal("EC2_IMAGE", imageRule.ResourceType)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableLockProtection", false)

	rbinRule := &RbinRule{
		svc:          mockClient,
		Identifier:   ptr.String("rule-abc123"),
		Description:  ptr.String("Retain EBS snapshots"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    "",
		settings:     settings,
	}

	mockClient.
		On("DeleteRule", mock.Anything, &rbin.DeleteRuleInput{
			Identifier: rbinRule.Identifier,
		}).
		Return(&rbin.DeleteRuleOutput{}, nil)

	err := rbinRule.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Remove_Locked_WithSetting(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableLockProtection", true)

	rbinRule := &RbinRule{
		svc:          mockClient,
		Identifier:   ptr.String("rule-locked-1"),
		Description:  ptr.String("Locked rule"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    string(rbintypes.LockStateLocked),
		settings:     settings,
	}

	lockEnd := time.Now().Add(7 * 24 * time.Hour)
	mockClient.
		On("UnlockRule", mock.Anything, &rbin.UnlockRuleInput{
			Identifier: rbinRule.Identifier,
		}).
		Return(&rbin.UnlockRuleOutput{
			LockState:   rbintypes.LockStatePendingUnlock,
			LockEndTime: &lockEnd,
		}, nil)

	err := rbinRule.Remove(context.TODO())
	assertions.Error(err)
	assertions.Contains(err.Error(), "unlock delay expires")
	assertions.Equal(string(rbintypes.LockStatePendingUnlock), rbinRule.LockState)
	assertions.NotNil(rbinRule.LockEndTime)

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Remove_Locked_SecondCallUsesPendingUnlock(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableLockProtection", true)

	rbinRule := &RbinRule{
		svc:          mockClient,
		Identifier:   ptr.String("rule-locked-loop"),
		Description:  ptr.String("Locked rule"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    string(rbintypes.LockStateLocked),
		settings:     settings,
	}

	lockEnd := time.Now().Add(7 * 24 * time.Hour)
	mockClient.
		On("UnlockRule", mock.Anything, &rbin.UnlockRuleInput{
			Identifier: rbinRule.Identifier,
		}).
		Return(&rbin.UnlockRuleOutput{
			LockState:   rbintypes.LockStatePendingUnlock,
			LockEndTime: &lockEnd,
		}, nil)

	// First call: unlocks the rule
	err := rbinRule.Remove(context.TODO())
	assertions.Error(err)
	assertions.Contains(err.Error(), "unlock delay expires")

	// Second call: should hit pending_unlock branch, NOT call UnlockRule again
	err = rbinRule.Remove(context.TODO())
	assertions.Error(err)
	assertions.Contains(err.Error(), "pending unlock")
	assertions.Contains(err.Error(), "day(s) remaining")

	// UnlockRule should have been called exactly once
	mockClient.AssertNumberOfCalls(t, "UnlockRule", 1)
	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Remove_Locked_WithoutSetting(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableLockProtection", false)

	rbinRule := &RbinRule{
		svc:          mockClient,
		Identifier:   ptr.String("rule-locked-2"),
		Description:  ptr.String("Locked rule"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    string(rbintypes.LockStateLocked),
		settings:     settings,
	}

	err := rbinRule.Remove(context.TODO())
	assertions.Error(err)
	assertions.Contains(err.Error(), "locked")
	assertions.Contains(err.Error(), "DisableLockProtection")

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Remove_PendingUnlock(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableLockProtection", true)

	lockEnd := time.Now().Add(5 * 24 * time.Hour)
	rbinRule := &RbinRule{
		svc:          mockClient,
		Identifier:   ptr.String("rule-pending-1"),
		Description:  ptr.String("Pending unlock rule"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    string(rbintypes.LockStatePendingUnlock),
		LockEndTime:  &lockEnd,
		settings:     settings,
	}

	err := rbinRule.Remove(context.TODO())
	assertions.Error(err)
	assertions.Contains(err.Error(), "pending unlock")
	assertions.Contains(err.Error(), "5 day(s) remaining")

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Remove_PendingUnlock_NoEndTime(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockRbinClient)

	settings := &libsettings.Setting{}
	settings.Set("DisableLockProtection", true)

	rbinRule := &RbinRule{
		svc:          mockClient,
		Identifier:   ptr.String("rule-pending-2"),
		Description:  ptr.String("Pending unlock rule"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    string(rbintypes.LockStatePendingUnlock),
		settings:     settings,
	}

	err := rbinRule.Remove(context.TODO())
	assertions.Error(err)
	assertions.Contains(err.Error(), "pending unlock")
	assertions.Contains(err.Error(), "unlock delay expires")

	mockClient.AssertExpectations(t)
}

func Test_Mock_RbinRule_Properties(t *testing.T) {
	assertions := assert.New(t)

	rbinRule := RbinRule{
		Identifier:   ptr.String("rule-abc123"),
		Description:  ptr.String("Retain EBS snapshots"),
		ResourceType: "EBS_SNAPSHOT",
		LockState:    "locked",
	}

	properties := rbinRule.Properties()

	assertions.Equal("rule-abc123", properties.Get("Identifier"))
	assertions.Equal("Retain EBS snapshots", properties.Get("Description"))
	assertions.Equal("EBS_SNAPSHOT", properties.Get("ResourceType"))
	assertions.Equal("locked", properties.Get("LockState"))
}

func Test_Mock_RbinRule_String(t *testing.T) {
	assertions := assert.New(t)

	rbinRule := RbinRule{
		Identifier: ptr.String("rule-abc123"),
	}

	assertions.Equal("rule-abc123", rbinRule.String())
}
