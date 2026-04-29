package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/omics"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OmicsWorkflowResource = "OmicsWorkflow"

func init() {
	registry.Register(&registry.Registration{
		Name:     OmicsWorkflowResource,
		Scope:    nuke.Account,
		Resource: &OmicsWorkflow{},
		Lister:   &OmicsWorkflowLister{},
	})
}

type OmicsWorkflowLister struct {
	svc OmicsClient
}

func (l *OmicsWorkflowLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = omics.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := omics.NewListWorkflowsPaginator(svc, &omics.ListWorkflowsInput{})
	for paginator.HasMorePages() {
		workflowsOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, workflow := range workflowsOutput.Items {
			resources = append(resources, &OmicsWorkflow{
				svc:  svc,
				ID:   workflow.Id,
				Name: workflow.Name,
			})
		}
	}

	return resources, nil
}

type OmicsWorkflow struct {
	svc  OmicsClient
	ID   *string
	Name *string
}

func (r *OmicsWorkflow) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkflow(ctx, &omics.DeleteWorkflowInput{
		Id: r.ID,
	})
	return err
}

func (r *OmicsWorkflow) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OmicsWorkflow) String() string {
	return *r.ID
}
