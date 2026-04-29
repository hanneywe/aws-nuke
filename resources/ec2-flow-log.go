package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EC2FlowLogResource = "EC2FlowLog"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2FlowLogResource,
		Scope:    nuke.Account,
		Resource: &EC2FlowLog{},
		Lister:   &EC2FlowLogLister{},
	})
}

type EC2FlowLogLister struct {
	svc EC2Client
}

func (l *EC2FlowLogLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeFlowLogsPaginator(svc,
		&ec2.DescribeFlowLogsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.FlowLogs {
			resources = append(resources, &EC2FlowLog{
				svc:                svc,
				FlowLogID:          resp.FlowLogs[i].FlowLogId,
				ResourceID:         resp.FlowLogs[i].ResourceId,
				TrafficType:        string(resp.FlowLogs[i].TrafficType),
				LogDestinationType: string(resp.FlowLogs[i].LogDestinationType),
				Tags:               resp.FlowLogs[i].Tags,
			})
		}
	}

	return resources, nil
}

type EC2FlowLog struct {
	svc                EC2Client
	FlowLogID          *string `property:"name=FlowLogId"`
	ResourceID         *string `property:"name=ResourceId"`
	TrafficType        string
	LogDestinationType string
	Tags               []ec2types.Tag
}

func (r *EC2FlowLog) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteFlowLogs(ctx, &ec2.DeleteFlowLogsInput{
		FlowLogIds: []string{*r.FlowLogID},
	})
	return err
}

func (r *EC2FlowLog) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2FlowLog) String() string {
	return *r.FlowLogID
}
