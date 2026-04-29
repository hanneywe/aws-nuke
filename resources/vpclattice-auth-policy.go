package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const VPCLatticeAuthPolicyResource = "VPCLatticeAuthPolicy"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeAuthPolicyResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeAuthPolicy{},
		Lister:   &VPCLatticeAuthPolicyLister{},
	})
}

type VPCLatticeAuthPolicyLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeAuthPolicyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	// Check service networks for auth policies
	snResources, err := l.listServiceNetworkAuthPolicies(ctx, svc)
	if err != nil {
		return nil, err
	}
	resources = append(resources, snResources...)

	// Check services for auth policies
	svcResources, err := l.listServiceAuthPolicies(ctx, svc)
	if err != nil {
		return nil, err
	}
	resources = append(resources, svcResources...)

	return resources, nil
}

func (l *VPCLatticeAuthPolicyLister) listServiceNetworkAuthPolicies(
	ctx context.Context, svc VPCLatticeClient,
) ([]resource.Resource, error) {
	var resources []resource.Resource

	snPaginator := vpclattice.NewListServiceNetworksPaginator(svc, &vpclattice.ListServiceNetworksInput{
		MaxResults: aws.Int32(100),
	})

	for snPaginator.HasMorePages() {
		resp, err := snPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, sn := range resp.Items {
			if sn.Arn == nil {
				continue
			}
			policyResp, err := svc.GetAuthPolicy(ctx, &vpclattice.GetAuthPolicyInput{
				ResourceIdentifier: sn.Arn,
			})
			if err != nil {
				continue // no policy attached
			}
			if policyResp.Policy == nil || *policyResp.Policy == "" {
				continue
			}
			resources = append(resources, &VPCLatticeAuthPolicy{
				svc:          svc,
				ResourceARN:  sn.Arn,
				ResourceType: aws.String("SERVICE_NETWORK"),
			})
		}
	}

	return resources, nil
}

func (l *VPCLatticeAuthPolicyLister) listServiceAuthPolicies(
	ctx context.Context, svc VPCLatticeClient,
) ([]resource.Resource, error) {
	var resources []resource.Resource

	svcPaginator := vpclattice.NewListServicesPaginator(svc, &vpclattice.ListServicesInput{
		MaxResults: aws.Int32(100),
	})

	for svcPaginator.HasMorePages() {
		resp, err := svcPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, service := range resp.Items {
			if service.Arn == nil {
				continue
			}
			policyResp, err := svc.GetAuthPolicy(ctx, &vpclattice.GetAuthPolicyInput{
				ResourceIdentifier: service.Arn,
			})
			if err != nil {
				continue // no policy attached
			}
			if policyResp.Policy == nil || *policyResp.Policy == "" {
				continue
			}
			resources = append(resources, &VPCLatticeAuthPolicy{
				svc:          svc,
				ResourceARN:  service.Arn,
				ResourceType: aws.String("SERVICE"),
			})
		}
	}

	return resources, nil
}

type VPCLatticeAuthPolicy struct {
	svc          VPCLatticeClient
	ResourceARN  *string
	ResourceType *string
}

func (r *VPCLatticeAuthPolicy) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAuthPolicy(ctx, &vpclattice.DeleteAuthPolicyInput{
		ResourceIdentifier: r.ResourceARN,
	})
	return err
}

func (r *VPCLatticeAuthPolicy) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeAuthPolicy) String() string {
	return *r.ResourceARN
}
