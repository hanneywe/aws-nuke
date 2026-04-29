package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/bedrockdataautomation"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BedrockDataAutomationProjectResource = "BedrockDataAutomationProject"

func init() {
	registry.Register(&registry.Registration{
		Name:     BedrockDataAutomationProjectResource,
		Scope:    nuke.Account,
		Resource: &BedrockDataAutomationProject{},
		Lister:   &BedrockDataAutomationProjectLister{},
	})
}

type BedrockDataAutomationProjectLister struct {
	svc BedrockDataAutomationClient
}

func (l *BedrockDataAutomationProjectLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = bedrockdataautomation.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := bedrockdataautomation.NewListDataAutomationProjectsPaginator(svc, &bedrockdataautomation.ListDataAutomationProjectsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range resp.Projects {
			resources = append(resources, &BedrockDataAutomationProject{
				svc:         svc,
				ProjectArn:  item.ProjectArn,
				ProjectName: item.ProjectName,
			})
		}
	}

	return resources, nil
}

type BedrockDataAutomationProject struct {
	svc         BedrockDataAutomationClient
	ProjectArn  *string
	ProjectName *string
}

func (r *BedrockDataAutomationProject) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteDataAutomationProject(ctx, &bedrockdataautomation.DeleteDataAutomationProjectInput{
		ProjectArn: r.ProjectArn,
	})
	return err
}

func (r *BedrockDataAutomationProject) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BedrockDataAutomationProject) String() string {
	return *r.ProjectArn
}
