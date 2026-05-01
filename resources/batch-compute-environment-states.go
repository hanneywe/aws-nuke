package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BatchComputeEnvironmentStateResource = "BatchComputeEnvironmentState"

func init() {
	registry.Register(&registry.Registration{
		Name:     BatchComputeEnvironmentStateResource,
		Scope:    nuke.Account,
		Resource: &BatchComputeEnvironmentState{},
		Lister:   &BatchComputeEnvironmentStateLister{},
	})
}

type BatchComputeEnvironmentStateLister struct {
	svc BatchComputeEnvironmentClient
}

func (l *BatchComputeEnvironmentStateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = batch.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := batch.NewDescribeComputeEnvironmentsPaginator(svc, &batch.DescribeComputeEnvironmentsInput{
		MaxResults: nil,
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range output.ComputeEnvironments {
			ce := &output.ComputeEnvironments[i]
			resources = append(resources, &BatchComputeEnvironmentState{
				svc:    svc,
				Name:   ce.ComputeEnvironmentName,
				State:  string(ce.State),
				Status: string(ce.Status),
				Type:   string(ce.Type),
				Tags:   ce.Tags,
			})
		}
	}

	return resources, nil
}

type BatchComputeEnvironmentState struct {
	svc    BatchComputeEnvironmentClient
	Name   *string
	State  string
	Status string
	Type   string
	Tags   map[string]string
}

func (r *BatchComputeEnvironmentState) Remove(ctx context.Context) error {
	_, err := r.svc.UpdateComputeEnvironment(ctx, &batch.UpdateComputeEnvironmentInput{
		ComputeEnvironment: r.Name,
		State:              batchtypes.CEStateDisabled,
	})
	return err
}

func (r *BatchComputeEnvironmentState) Filter() error {
	if r.State == string(batchtypes.CEStateDisabled) {
		return fmt.Errorf("already disabled")
	}
	if r.Status == string(batchtypes.CEStatusDeleting) || r.Status == string(batchtypes.CEStatusDeleted) {
		return fmt.Errorf("compute environment is being deleted")
	}
	return nil
}

func (r *BatchComputeEnvironmentState) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BatchComputeEnvironmentState) String() string {
	return *r.Name
}
