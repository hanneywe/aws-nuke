package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeListenerResource = "VPCLatticeListener"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeListenerResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeListener{},
		Lister:   &VPCLatticeListenerLister{},
		DependsOn: []string{
			VPCLatticeListenerRuleResource,
		},
	})
}

type VPCLatticeListenerLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeListenerLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	servicesPaginator := vpclattice.NewListServicesPaginator(svc, &vpclattice.ListServicesInput{
		MaxResults: aws.Int32(100),
	})

	for servicesPaginator.HasMorePages() {
		servicesResp, err := servicesPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, service := range servicesResp.Items {
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
					var tags map[string]string
					if listener.Arn != nil {
						tagsResp, err := svc.ListTagsForResource(ctx, &vpclattice.ListTagsForResourceInput{
							ResourceArn: listener.Arn,
						})
						if err != nil {
							opts.Logger.Warnf("unable to fetch tags for listener: %s", *listener.Arn)
						} else {
							tags = tagsResp.Tags
						}
					}

					resources = append(resources, &VPCLatticeListener{
						svc:         svc,
						ID:          listener.Id,
						ARN:         listener.Arn,
						Name:        listener.Name,
						ServiceID:   service.Id,
						ServiceName: service.Name,
						Port:        listener.Port,
						Protocol:    aws.String(string(listener.Protocol)),
						Tags:        tags,
					})
				}
			}
		}
	}

	return resources, nil
}

type VPCLatticeListener struct {
	svc         VPCLatticeClient
	ID          *string
	ARN         *string
	Name        *string
	ServiceID   *string
	ServiceName *string
	Port        *int32
	Protocol    *string
	Tags        map[string]string
}

func (r *VPCLatticeListener) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteListener(ctx, &vpclattice.DeleteListenerInput{
		ServiceIdentifier:  r.ServiceID,
		ListenerIdentifier: r.ID,
	})
	return err
}

func (r *VPCLatticeListener) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeListener) String() string {
	return fmt.Sprintf("%s -> %s", *r.ServiceName, *r.Name)
}
