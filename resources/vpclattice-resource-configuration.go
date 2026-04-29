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

const VPCLatticeResourceConfigurationResource = "VPCLatticeResourceConfiguration"

func init() {
	registry.Register(&registry.Registration{
		Name:     VPCLatticeResourceConfigurationResource,
		Scope:    nuke.Account,
		Resource: &VPCLatticeResourceConfiguration{},
		Lister:   &VPCLatticeResourceConfigurationLister{},
	})
}

type VPCLatticeResourceConfigurationLister struct {
	svc VPCLatticeClient
}

func (l *VPCLatticeResourceConfigurationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = vpclattice.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := vpclattice.NewListResourceConfigurationsPaginator(svc, &vpclattice.ListResourceConfigurationsInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Items {
			resources = append(resources, &VPCLatticeResourceConfiguration{
				svc:  svc,
				ID:   resp.Items[i].Id,
				ARN:  resp.Items[i].Arn,
				Name: resp.Items[i].Name,
				Type: aws.String(string(resp.Items[i].Type)),
			})
		}
	}

	return resources, nil
}

type VPCLatticeResourceConfiguration struct {
	svc  VPCLatticeClient
	ID   *string
	ARN  *string
	Name *string
	Type *string
}

func (r *VPCLatticeResourceConfiguration) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteResourceConfiguration(ctx, &vpclattice.DeleteResourceConfigurationInput{
		ResourceConfigurationIdentifier: r.ARN,
	})
	return err
}

func (r *VPCLatticeResourceConfiguration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *VPCLatticeResourceConfiguration) String() string {
	return *r.Name
}
