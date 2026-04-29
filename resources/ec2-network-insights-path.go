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

const EC2NetworkInsightsPathResource = "EC2NetworkInsightsPath"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2NetworkInsightsPathResource,
		Scope:    nuke.Account,
		Resource: &EC2NetworkInsightsPath{},
		Lister:   &EC2NetworkInsightsPathLister{},
		DependsOn: []string{
			EC2NetworkInsightsAnalysisResource,
		},
	})
}

type EC2NetworkInsightsPathLister struct {
	svc EC2Client
}

func (l *EC2NetworkInsightsPathLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeNetworkInsightsPathsPaginator(svc,
		&ec2.DescribeNetworkInsightsPathsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.NetworkInsightsPaths {
			resources = append(resources, &EC2NetworkInsightsPath{
				svc:                   svc,
				NetworkInsightsPathID: resp.NetworkInsightsPaths[i].NetworkInsightsPathId,
				Source:                resp.NetworkInsightsPaths[i].Source,
				Destination:           resp.NetworkInsightsPaths[i].Destination,
				Protocol:              string(resp.NetworkInsightsPaths[i].Protocol),
				Tags:                  resp.NetworkInsightsPaths[i].Tags,
			})
		}
	}

	return resources, nil
}

type EC2NetworkInsightsPath struct {
	svc                   EC2Client
	NetworkInsightsPathID *string `property:"name=NetworkInsightsPathId"`
	Source                *string
	Destination           *string
	Protocol              string
	Tags                  []ec2types.Tag
}

func (r *EC2NetworkInsightsPath) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteNetworkInsightsPath(ctx, &ec2.DeleteNetworkInsightsPathInput{
		NetworkInsightsPathId: r.NetworkInsightsPathID,
	})
	return err
}

func (r *EC2NetworkInsightsPath) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2NetworkInsightsPath) String() string {
	return *r.NetworkInsightsPathID
}
