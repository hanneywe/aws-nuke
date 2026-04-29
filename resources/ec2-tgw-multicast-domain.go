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

const EC2TGWMulticastDomainResource = "EC2TGWMulticastDomain"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TGWMulticastDomainResource,
		Scope:    nuke.Account,
		Resource: &EC2TGWMulticastDomain{},
		Lister:   &EC2TGWMulticastDomainLister{},
	})
}

type EC2TGWMulticastDomainLister struct {
	svc EC2Client
}

func (l *EC2TGWMulticastDomainLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeTransitGatewayMulticastDomainsPaginator(svc,
		&ec2.DescribeTransitGatewayMulticastDomainsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, domain := range resp.TransitGatewayMulticastDomains {
			resources = append(resources, &EC2TGWMulticastDomain{
				svc:                             svc,
				TransitGatewayMulticastDomainID: domain.TransitGatewayMulticastDomainId,
				TransitGatewayID:                domain.TransitGatewayId,
				State:                           string(domain.State),
				Tags:                            domain.Tags,
			})
		}
	}

	return resources, nil
}

type EC2TGWMulticastDomain struct {
	svc                             EC2Client
	TransitGatewayMulticastDomainID *string `property:"name=TransitGatewayMulticastDomainId"`
	TransitGatewayID                *string `property:"name=TransitGatewayId"`
	State                           string
	Tags                            []ec2types.Tag
}

func (r *EC2TGWMulticastDomain) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTransitGatewayMulticastDomain(ctx, &ec2.DeleteTransitGatewayMulticastDomainInput{
		TransitGatewayMulticastDomainId: r.TransitGatewayMulticastDomainID,
	})
	return err
}

func (r *EC2TGWMulticastDomain) Filter() error {
	if r.State == string(ec2types.TransitGatewayMulticastDomainStateDeleted) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2TGWMulticastDomain) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TGWMulticastDomain) String() string {
	return *r.TransitGatewayMulticastDomainID
}
