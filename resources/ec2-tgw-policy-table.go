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

const EC2TGWPolicyTableResource = "EC2TGWPolicyTable"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2TGWPolicyTableResource,
		Scope:    nuke.Account,
		Resource: &EC2TGWPolicyTable{},
		Lister:   &EC2TGWPolicyTableLister{},
	})
}

type EC2TGWPolicyTableLister struct {
	svc EC2Client
}

func (l *EC2TGWPolicyTableLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeTransitGatewayPolicyTablesPaginator(svc,
		&ec2.DescribeTransitGatewayPolicyTablesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, policyTable := range resp.TransitGatewayPolicyTables {
			resources = append(resources, &EC2TGWPolicyTable{
				svc:                         svc,
				TransitGatewayPolicyTableID: policyTable.TransitGatewayPolicyTableId,
				TransitGatewayID:            policyTable.TransitGatewayId,
				State:                       string(policyTable.State),
				Tags:                        policyTable.Tags,
			})
		}
	}

	return resources, nil
}

type EC2TGWPolicyTable struct {
	svc                         EC2Client
	TransitGatewayPolicyTableID *string `property:"name=TransitGatewayPolicyTableId"`
	TransitGatewayID            *string `property:"name=TransitGatewayId"`
	State                       string
	Tags                        []ec2types.Tag
}

func (r *EC2TGWPolicyTable) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteTransitGatewayPolicyTable(ctx, &ec2.DeleteTransitGatewayPolicyTableInput{
		TransitGatewayPolicyTableId: r.TransitGatewayPolicyTableID,
	})
	return err
}

func (r *EC2TGWPolicyTable) Filter() error {
	if r.State == string(ec2types.TransitGatewayPolicyTableStateDeleted) {
		return fmt.Errorf("already deleted")
	}
	return nil
}

func (r *EC2TGWPolicyTable) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2TGWPolicyTable) String() string {
	return *r.TransitGatewayPolicyTableID
}
