package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const Route53RecoveryReadinessRecoveryGroupResource = "Route53RecoveryReadinessRecoveryGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     Route53RecoveryReadinessRecoveryGroupResource,
		Scope:    nuke.Account,
		Resource: &Route53RecoveryReadinessRecoveryGroup{},
		Lister:   &Route53RecoveryReadinessRecoveryGroupLister{},
	})
}

type Route53RecoveryReadinessRecoveryGroupLister struct {
	svc Route53RecoveryReadinessClient
}

func (l *Route53RecoveryReadinessRecoveryGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = route53recoveryreadiness.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &route53recoveryreadiness.ListRecoveryGroupsInput{}
	for {
		resp, err := svc.ListRecoveryGroups(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, rg := range resp.RecoveryGroups {
			resources = append(resources, &Route53RecoveryReadinessRecoveryGroup{
				svc:               svc,
				RecoveryGroupName: rg.RecoveryGroupName,
				RecoveryGroupArn:  rg.RecoveryGroupArn,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type Route53RecoveryReadinessRecoveryGroup struct {
	svc               Route53RecoveryReadinessClient
	RecoveryGroupName *string
	RecoveryGroupArn  *string
}

func (r *Route53RecoveryReadinessRecoveryGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRecoveryGroup(ctx, &route53recoveryreadiness.DeleteRecoveryGroupInput{
		RecoveryGroupName: r.RecoveryGroupName,
	})
	return err
}

func (r *Route53RecoveryReadinessRecoveryGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *Route53RecoveryReadinessRecoveryGroup) String() string {
	return *r.RecoveryGroupName
}
