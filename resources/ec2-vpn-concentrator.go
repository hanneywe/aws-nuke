package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2VPNConcentratorResource = "EC2VPNConcentrator"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2VPNConcentratorResource,
		Scope:    nuke.Account,
		Resource: &EC2VPNConcentrator{},
		Lister:   &EC2VPNConcentratorLister{},
		DependsOn: []string{
			EC2VPNConnectionResource,
		},
	})
}

type EC2VPNConcentratorLister struct {
	svc EC2Client
}

func (l *EC2VPNConcentratorLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeVpnConcentratorsPaginator(svc, &ec2.DescribeVpnConcentratorsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.VpnConcentrators {
			vc := &resp.VpnConcentrators[i]
			resources = append(resources, &EC2VPNConcentrator{
				svc:               svc,
				VpnConcentratorID: vc.VpnConcentratorId,
				TransitGatewayID:  vc.TransitGatewayId,
				State:             vc.State,
				Tags:              vc.Tags,
			})
		}
	}

	return resources, nil
}

type EC2VPNConcentrator struct {
	svc               EC2Client
	VpnConcentratorID *string `property:"name=VpnConcentratorId"`
	TransitGatewayID  *string `property:"name=TransitGatewayId"`
	State             *string
	Tags              []ec2types.Tag
}

func (r *EC2VPNConcentrator) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteVpnConcentrator(ctx, &ec2.DeleteVpnConcentratorInput{
		VpnConcentratorId: r.VpnConcentratorID,
	})
	return err
}

func (r *EC2VPNConcentrator) Filter() error {
	if r.State != nil && *r.State == "deleted" {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2VPNConcentrator) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2VPNConcentrator) String() string {
	return *r.VpnConcentratorID
}
