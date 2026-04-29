package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const TransferWorkflowResource = "TransferWorkflow"

func init() {
	registry.Register(&registry.Registration{
		Name:     TransferWorkflowResource,
		Scope:    nuke.Account,
		Resource: &TransferWorkflow{},
		Lister:   &TransferWorkflowLister{},
	})
}

type TransferWorkflowLister struct {
	svc TransferClient
}

func (l *TransferWorkflowLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = transfer.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := transfer.NewListWorkflowsPaginator(svc, &transfer.ListWorkflowsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.Workflows {
			item := &resp.Workflows[i]
			resources = append(resources, &TransferWorkflow{
				svc:         svc,
				WorkflowID:  item.WorkflowId,
				Description: item.Description,
				Arn:         item.Arn,
			})
		}
	}

	return resources, nil
}

type TransferWorkflow struct {
	svc         TransferClient
	WorkflowID  *string
	Description *string
	Arn         *string
}

func (r *TransferWorkflow) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWorkflow(ctx, &transfer.DeleteWorkflowInput{
		WorkflowId: r.WorkflowID,
	})
	return err
}

func (r *TransferWorkflow) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *TransferWorkflow) String() string {
	return *r.WorkflowID
}
