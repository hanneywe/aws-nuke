package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/rds"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const RDSDBSecurityGroupResource = "RDSDBSecurityGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     RDSDBSecurityGroupResource,
		Scope:    nuke.Account,
		Resource: &RDSDBSecurityGroup{},
		Lister:   &RDSDBSecurityGroupLister{},
	})
}

type RDSDBSecurityGroupLister struct {
	svc RDSV2Client
}

func (l *RDSDBSecurityGroupLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = rds.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := rds.NewDescribeDBSecurityGroupsPaginator(svc, &rds.DescribeDBSecurityGroupsInput{})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, securityGroup := range output.DBSecurityGroups {
			resources = append(resources, &RDSDBSecurityGroup{
				svc:                 svc,
				DBSecurityGroupName: securityGroup.DBSecurityGroupName,
				DBSecurityGroupArn:  securityGroup.DBSecurityGroupArn,
			})
		}
	}

	return resources, nil
}

type RDSDBSecurityGroup struct {
	svc                 RDSV2Client
	DBSecurityGroupName *string
	DBSecurityGroupArn  *string
}

func (r *RDSDBSecurityGroup) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDBSecurityGroup(ctx, &rds.DeleteDBSecurityGroupInput{
		DBSecurityGroupName: r.DBSecurityGroupName,
	})
	return err
}

func (r *RDSDBSecurityGroup) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *RDSDBSecurityGroup) String() string {
	return *r.DBSecurityGroupName
}

func (r *RDSDBSecurityGroup) Filter() error {
	if r.DBSecurityGroupName != nil && *r.DBSecurityGroupName == "default" {
		return fmt.Errorf("cannot delete default DB security group")
	}
	return nil
}
