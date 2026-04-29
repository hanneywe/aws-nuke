package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/appstream"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AppStreamEntitlementResource = "AppStreamEntitlement"

func init() {
	registry.Register(&registry.Registration{
		Name:     AppStreamEntitlementResource,
		Scope:    nuke.Account,
		Resource: &AppStreamEntitlement{},
		Lister:   &AppStreamEntitlementLister{},
	})
}

type AppStreamEntitlementLister struct {
	svc AppStreamClient
}

func (l *AppStreamEntitlementLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = appstream.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	stackResp, err := svc.DescribeStacks(ctx, &appstream.DescribeStacksInput{})
	if err != nil {
		return nil, err
	}
	for i := range stackResp.Stacks {
		stack := &stackResp.Stacks[i]
		entitlementResp, err := svc.DescribeEntitlements(ctx, &appstream.DescribeEntitlementsInput{
			StackName: stack.Name,
		})
		if err != nil {
			return nil, err
		}
		for _, entitlement := range entitlementResp.Entitlements {
			resources = append(resources, &AppStreamEntitlement{
				svc:         svc,
				StackName:   stack.Name,
				Name:        entitlement.Name,
				Description: entitlement.Description,
			})
		}
	}

	return resources, nil
}

type AppStreamEntitlement struct {
	svc         AppStreamClient
	StackName   *string
	Name        *string
	Description *string
}

func (r *AppStreamEntitlement) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEntitlement(ctx, &appstream.DeleteEntitlementInput{
		Name:      r.Name,
		StackName: r.StackName,
	})
	return err
}

func (r *AppStreamEntitlement) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AppStreamEntitlement) String() string {
	return *r.Name
}
