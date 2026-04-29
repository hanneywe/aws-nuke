package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/gotidy/ptr"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeListenerRuleResource = "VPCLatticeListenerRule"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeListenerRuleResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeListenerRule{},
		Lister:   &VPCLatticeListenerRuleLister{},
	})
}

type VPCLatticeListenerRuleLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeListenerRuleLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	// First, list all services
	servicesPaginator := vpclattice.NewListServicesPaginator(svc, &vpclattice.ListServicesInput{
		MaxResults: aws.Int32(100),
	})

	for servicesPaginator.HasMorePages() {
		servicesResp, err := servicesPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, service := range servicesResp.Items {
			// For each service, list all listeners
			listenersPaginator := vpclattice.NewListListenersPaginator(svc, &vpclattice.ListListenersInput{
				ServiceIdentifier: service.Id,
				MaxResults:        aws.Int32(100),
			})

			for listenersPaginator.HasMorePages() {
				listenersResp, err := listenersPaginator.NextPage(ctx)
				if err != nil {
					opts.Logger.Warnf("unable to list listeners for service %s: %v", *service.Id, err)
					break
				}

				for _, listener := range listenersResp.Items {
					// For each listener, list all rules
					rulesPaginator := vpclattice.NewListRulesPaginator(svc, &vpclattice.ListRulesInput{
						ServiceIdentifier:  service.Id,
						ListenerIdentifier: listener.Id,
						MaxResults:         aws.Int32(100),
					})

					for rulesPaginator.HasMorePages() {
						rulesResp, err := rulesPaginator.NextPage(ctx)
						if err != nil {
							opts.Logger.Warnf("unable to list rules for listener %s on service %s: %v", *listener.Id, *service.Id, err)
							break
						}

						for _, rule := range rulesResp.Items {
							resources = append(resources, &VPCLatticeListenerRule{
								svc:          svc,
								ID:           rule.Id,
								ARN:          rule.Arn,
								Name:         rule.Name,
								ServiceID:    service.Id,
								ServiceName:  service.Name,
								ListenerID:   listener.Id,
								ListenerName: listener.Name,
								Priority:     rule.Priority,
								IsDefault:    ptr.ToBool(rule.IsDefault),
							})
						}
					}
				}
			}
		}
	}

	return resources, nil
}

type VPCLatticeListenerRule struct {
	svc          VPCLatticeClient
	ID           *string
	ARN          *string
	Name         *string
	ServiceID    *string
	ServiceName  *string
	ListenerID   *string
	ListenerName *string
	Priority     *int32
	IsDefault    bool
	Tags         map[string]string
}

func (r *VPCLatticeListenerRule) Filter() error {
	if r.IsDefault {
		return fmt.Errorf("cannot delete default listener rule")
	}
	return nil
}

func (r *VPCLatticeListenerRule) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRule(ctx, &vpclattice.DeleteRuleInput{
		ServiceIdentifier:  r.ServiceID,
		ListenerIdentifier: r.ListenerID,
		RuleIdentifier:     r.ID,
	})
	return err
}

func (r *VPCLatticeListenerRule) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeListenerRule) String() string {
	return fmt.Sprintf("%s -> %s -> %s", *r.ServiceName, *r.ListenerName, *r.Name)
}
