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

const EC2TGWPrefixListReferenceResource = "EC2TGWPrefixListReference"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TGWPrefixListReferenceResource,
		Scope:    nuke.Account,
		Resource: &EC2TGWPrefixListReference{},
		Lister:   &EC2TGWPrefixListReferenceLister{},
	})
}

type EC2TGWPrefixListReferenceLister struct {
	svc EC2Client
}

func (l *EC2TGWPrefixListReferenceLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First, enumerate all TGW route tables
	routeTablePaginator := ec2.NewDescribeTransitGatewayRouteTablesPaginator(svc,
		&ec2.DescribeTransitGatewayRouteTablesInput{})

	for routeTablePaginator.HasMorePages() {
		routeTablePage, err := routeTablePaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, routeTable := range routeTablePage.TransitGatewayRouteTables {
			// For each route table, get its prefix list references
			prefixListPaginator := ec2.NewGetTransitGatewayPrefixListReferencesPaginator(svc,
				&ec2.GetTransitGatewayPrefixListReferencesInput{
					TransitGatewayRouteTableId: routeTable.TransitGatewayRouteTableId,
				})

			for prefixListPaginator.HasMorePages() {
				prefixListPage, err := prefixListPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, ref := range prefixListPage.TransitGatewayPrefixListReferences {
					resources = append(resources, &EC2TGWPrefixListReference{
						svc:                        svc,
						TransitGatewayRouteTableID: routeTable.TransitGatewayRouteTableId,
						PrefixListID:               ref.PrefixListId,
						State:                      string(ref.State),
					})
				}
			}
		}
	}

	return resources, nil
}

type EC2TGWPrefixListReference struct {
	svc                        EC2Client
	TransitGatewayRouteTableID *string `property:"name=TransitGatewayRouteTableId"`
	PrefixListID               *string `property:"name=PrefixListId"`
	State                      string
}

func (r *EC2TGWPrefixListReference) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTransitGatewayPrefixListReference(ctx, &ec2.DeleteTransitGatewayPrefixListReferenceInput{
		TransitGatewayRouteTableId: r.TransitGatewayRouteTableID,
		PrefixListId:               r.PrefixListID,
	})
	return err
}

func (r *EC2TGWPrefixListReference) Filter() error {
	if r.State == string(ec2types.TransitGatewayPrefixListReferenceStateDeleting) {
		return fmt.Errorf("already being deleted")
	}
	return nil
}

func (r *EC2TGWPrefixListReference) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TGWPrefixListReference) String() string {
	return fmt.Sprintf("%s -> %s", *r.TransitGatewayRouteTableID, *r.PrefixListID)
}
