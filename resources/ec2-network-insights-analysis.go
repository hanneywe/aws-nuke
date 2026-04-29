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

const EC2NetworkInsightsAnalysisResource = "EC2NetworkInsightsAnalysis"

func init() {
	registry.Register(&registry.Registration{
		Name:     EC2NetworkInsightsAnalysisResource,
		Scope:    nuke.Account,
		Resource: &EC2NetworkInsightsAnalysis{},
		Lister:   &EC2NetworkInsightsAnalysisLister{},
	})
}

type EC2NetworkInsightsAnalysisLister struct {
	svc EC2Client
}

func (l *EC2NetworkInsightsAnalysisLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ec2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := ec2.NewDescribeNetworkInsightsAnalysesPaginator(svc,
		&ec2.DescribeNetworkInsightsAnalysesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.NetworkInsightsAnalyses {
			resources = append(resources, &EC2NetworkInsightsAnalysis{
				svc:                       svc,
				NetworkInsightsAnalysisID: resp.NetworkInsightsAnalyses[i].NetworkInsightsAnalysisId,
				NetworkInsightsPathID:     resp.NetworkInsightsAnalyses[i].NetworkInsightsPathId,
				Status:                    string(resp.NetworkInsightsAnalyses[i].Status),
				Tags:                      resp.NetworkInsightsAnalyses[i].Tags,
			})
		}
	}

	return resources, nil
}

type EC2NetworkInsightsAnalysis struct {
	svc                       EC2Client
	NetworkInsightsAnalysisID *string `property:"name=NetworkInsightsAnalysisId"`
	NetworkInsightsPathID     *string `property:"name=NetworkInsightsPathId"`
	Status                    string
	Tags                      []ec2types.Tag
}

func (r *EC2NetworkInsightsAnalysis) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteNetworkInsightsAnalysis(ctx, &ec2.DeleteNetworkInsightsAnalysisInput{
		NetworkInsightsAnalysisId: r.NetworkInsightsAnalysisID,
	})
	return err
}

func (r *EC2NetworkInsightsAnalysis) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EC2NetworkInsightsAnalysis) String() string {
	return *r.NetworkInsightsAnalysisID
}
