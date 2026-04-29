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

const EC2CarrierGatewayResource = "EC2CarrierGateway"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2CarrierGatewayResource,
		Scope:    nuke.Account,
		Resource: &EC2CarrierGateway{},
		Lister:   &EC2CarrierGatewayLister{},
	})
}

type EC2CarrierGatewayLister struct {
	svc EC2Client
}

func (l *EC2CarrierGatewayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeCarrierGatewaysPaginator(svc,
		&ec2.DescribeCarrierGatewaysInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, gateway := range resp.CarrierGateways {
			resources = append(resources, &EC2CarrierGateway{
				svc:              svc,
				CarrierGatewayID: gateway.CarrierGatewayId,
				VpcID:            gateway.VpcId,
				State:            string(gateway.State),
				OwnerID:          gateway.OwnerId,
				Tags:             gateway.Tags,
			})
		}
	}

	return resources, nil
}

type EC2CarrierGateway struct {
	svc              EC2Client
	CarrierGatewayID *string `property:"name=CarrierGatewayId"`
	VpcID            *string `property:"name=VpcId"`
	State            string
	OwnerID          *string `property:"name=OwnerId"`
	Tags             []ec2types.Tag
}

func (r *EC2CarrierGateway) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCarrierGateway(ctx, &ec2.DeleteCarrierGatewayInput{
		CarrierGatewayId: r.CarrierGatewayID,
	})
	return err
}

func (r *EC2CarrierGateway) Filter() error {
	if r.State == string(ec2types.CarrierGatewayStateDeleted) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2CarrierGateway) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2CarrierGateway) String() string {
	return *r.CarrierGatewayID
}
