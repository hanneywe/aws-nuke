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

const EC2TGWRouteTableResource = "EC2TGWRouteTable"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TGWRouteTableResource,
		Scope:    nuke.Account,
		Resource: &EC2TGWRouteTable{},
		Lister:   &EC2TGWRouteTableLister{},
		DependsOn: []string{
			EC2TGWPrefixListReferenceResource,
			EC2TGWAttachmentResource,
		},
	})
}

type EC2TGWRouteTableLister struct {
	svc EC2Client
}

func (l *EC2TGWRouteTableLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeTransitGatewayRouteTablesPaginator(svc,
		&ec2.DescribeTransitGatewayRouteTablesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, routeTable := range resp.TransitGatewayRouteTables {
			resources = append(resources, &EC2TGWRouteTable{
				svc:                        svc,
				TransitGatewayRouteTableID: routeTable.TransitGatewayRouteTableId,
				TransitGatewayID:           routeTable.TransitGatewayId,
				State:                      string(routeTable.State),
				DefaultAssociation:         routeTable.DefaultAssociationRouteTable,
				DefaultPropagation:         routeTable.DefaultPropagationRouteTable,
				Tags:                       routeTable.Tags,
			})
		}
	}

	return resources, nil
}

type EC2TGWRouteTable struct {
	svc                        EC2Client
	TransitGatewayRouteTableID *string `property:"name=TransitGatewayRouteTableId"`
	TransitGatewayID           *string `property:"name=TransitGatewayId"`
	State                      string
	DefaultAssociation         *bool
	DefaultPropagation         *bool
	Tags                       []ec2types.Tag
}

func (r *EC2TGWRouteTable) Remove(ctx context.Context) error {
	if (r.DefaultAssociation != nil && *r.DefaultAssociation) ||
		(r.DefaultPropagation != nil && *r.DefaultPropagation) {
		_, err := r.svc.ModifyTransitGateway(ctx, &ec2.ModifyTransitGatewayInput{
			TransitGatewayId: r.TransitGatewayID,
			Options: &ec2types.ModifyTransitGatewayOptions{
				DefaultRouteTableAssociation: ec2types.DefaultRouteTableAssociationValueDisable,
				DefaultRouteTablePropagation: ec2types.DefaultRouteTablePropagationValueDisable,
			},
		})
		if err != nil {
			return err
		}
		r.DefaultAssociation = nil
		r.DefaultPropagation = nil
	}

	_, err := r.svc.DeleteTransitGatewayRouteTable(ctx, &ec2.DeleteTransitGatewayRouteTableInput{
		TransitGatewayRouteTableId: r.TransitGatewayRouteTableID,
	})
	return err
}

func (r *EC2TGWRouteTable) Filter() error {
	if r.State == string(ec2types.TransitGatewayRouteTableStateDeleted) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2TGWRouteTable) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TGWRouteTable) String() string {
	return *r.TransitGatewayRouteTableID
}
