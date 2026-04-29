package resources

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/rbin"
	rbintypes "github.com/aws/aws-sdk-go-v2/service/rbin/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const RbinRuleResource = "RbinRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     RbinRuleResource,
		Scope:    nuke.Account,
		Resource: &RbinRule{},
		Lister:   &RbinRuleLister{},
		Settings: []string{
			"DisableLockProtection",
		},
	})
}

// rbinResourceTypes defines the known resource types to iterate over when listing rules.
// ListRules requires a ResourceType parameter, so we must call it once per type.
var rbinResourceTypes = []rbintypes.ResourceType{
	rbintypes.ResourceTypeEbsSnapshot,
	rbintypes.ResourceTypeEc2Image,
}

type RbinRuleLister struct {
	svc RbinClient
}

func (l *RbinRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = rbin.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	for _, resourceType := range rbinResourceTypes {
		paginator := rbin.NewListRulesPaginator(svc, &rbin.ListRulesInput{
			ResourceType: resourceType,
		})

		for paginator.HasMorePages() {
			listRulesOutput, err := paginator.NextPage(ctx)
			if err != nil {
				return nil, err
			}

			for _, ruleSummary := range listRulesOutput.Rules {
				r := &RbinRule{
					svc:          svc,
					Identifier:   ruleSummary.Identifier,
					Description:  ruleSummary.Description,
					ResourceType: string(resourceType),
					LockState:    string(ruleSummary.LockState),
				}

				// Fetch LockEndTime for rules that are pending unlock
				if ruleSummary.LockState == rbintypes.LockStatePendingUnlock {
					ruleDetail, err := svc.GetRule(ctx, &rbin.GetRuleInput{
						Identifier: ruleSummary.Identifier,
					})
					if err == nil {
						r.LockEndTime = ruleDetail.LockEndTime
					}
				}

				resources = append(resources, r)
			}
		}
	}

	return resources, nil
}

type RbinRule struct {
	svc          RbinClient
	Identifier   *string
	Description  *string
	ResourceType string
	LockState    string
	LockEndTime  *time.Time `property:"-"`
	settings     *libsettings.Setting
}

func (r *RbinRule) Remove(ctx context.Context) error {
	if r.LockState == string(rbintypes.LockStateLocked) {
		if !r.settings.GetBool("DisableLockProtection") {
			return fmt.Errorf("rule %s is locked, set DisableLockProtection to true to unlock and delete", *r.Identifier)
		}

		resp, err := r.svc.UnlockRule(ctx, &rbin.UnlockRuleInput{
			Identifier: r.Identifier,
		})
		if err != nil {
			return err
		}

		// Update local state so subsequent loop iterations see the correct lock state
		r.LockState = string(resp.LockState)
		r.LockEndTime = resp.LockEndTime

		// After unlocking, the rule enters pending_unlock state. The minimum
		// unlock delay is 7 days, so there is no point retrying in this run.
		return fmt.Errorf("rule %s unlocked, but cannot be deleted until the unlock delay expires (minimum 7 days)", *r.Identifier)
	}

	if r.LockState == string(rbintypes.LockStatePendingUnlock) {
		if r.LockEndTime != nil {
			remaining := int(math.Ceil(time.Until(*r.LockEndTime).Hours() / 24))
			return fmt.Errorf("rule %s is pending unlock, %d day(s) remaining until %s",
				*r.Identifier, remaining, r.LockEndTime.Format(time.RFC3339))
		}
		return fmt.Errorf("rule %s is pending unlock, cannot be deleted until the unlock delay expires", *r.Identifier)
	}

	_, err := r.svc.DeleteRule(ctx, &rbin.DeleteRuleInput{
		Identifier: r.Identifier,
	})
	return err
}

func (r *RbinRule) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RbinRule) String() string {
	return *r.Identifier
}

func (r *RbinRule) Settings(setting *libsettings.Setting) {
	r.settings = setting
}
