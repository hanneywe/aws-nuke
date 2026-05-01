package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/batch"
	batchtypes "github.com/aws/aws-sdk-go-v2/service/batch/types"

	liberror "github.com/ekristen/libnuke/pkg/errors"
	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BatchComputeEnvironmentResource = "BatchComputeEnvironment"

func init() {
	registry.Register(&registry.Registration{
		Name:     BatchComputeEnvironmentResource,
		Scope:    nuke.Account,
		Resource: &BatchComputeEnvironment{},
		Lister:   &BatchComputeEnvironmentLister{},
	})
}

type BatchComputeEnvironmentLister struct {
	svc BatchComputeEnvironmentClient
}

func (l *BatchComputeEnvironmentLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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
			if ce.Status == BatchCEStatusDeleted {
				continue
			}

			resources = append(resources, &BatchComputeEnvironment{
				svc:    svc,
				Name:   ce.ComputeEnvironmentName,
				Status: string(ce.Status),
				State:  string(ce.State),
				Type:   string(ce.Type),
				Tags:   ce.Tags,
			})
		}
	}

	return resources, nil
}

type BatchComputeEnvironment struct {
	svc    BatchComputeEnvironmentClient
	Name   *string
	Status string
	State  string
	Type   string
	Tags   map[string]string
}

func (r *BatchComputeEnvironment) Filter() error {
	if r.Status == string(batchtypes.CEStatusDeleting) {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *BatchComputeEnvironment) Remove(ctx context.Context) error {
	if r.State == string(batchtypes.CEStateEnabled) {
		_, err := r.svc.UpdateComputeEnvironment(ctx, &batch.UpdateComputeEnvironmentInput{
			ComputeEnvironment: r.Name,
			State:              batchtypes.CEStateDisabled,
		})
		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteComputeEnvironment(ctx, &batch.DeleteComputeEnvironmentInput{
		ComputeEnvironment: r.Name,
	})
	return err
}

func (r *BatchComputeEnvironment) HandleWait(ctx context.Context) error {
	resp, err := r.svc.DescribeComputeEnvironments(ctx, &batch.DescribeComputeEnvironmentsInput{
		ComputeEnvironments: []string{*r.Name},
	})
	if err != nil {
		return err
	}

	if len(resp.ComputeEnvironments) == 0 {
		return nil
	}

	ce := resp.ComputeEnvironments[0]

	switch ce.Status {
	case BatchCEStatusDeleted:
		return nil
	case batchtypes.CEStatusDeleting:
		return liberror.ErrWaitResource("waiting for compute environment to delete")
	case batchtypes.CEStatusInvalid:
		r.svc.DeleteComputeEnvironment(ctx, &batch.DeleteComputeEnvironmentInput{ //nolint:errcheck
			ComputeEnvironment: r.Name,
		})
		return liberror.ErrWaitResource("retrying delete of INVALID compute environment")
	default:
		return liberror.ErrWaitResource(fmt.Sprintf("waiting for compute environment to transition, current status: %s", ce.Status))
	}
}

func (r *BatchComputeEnvironment) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BatchComputeEnvironment) String() string {
	return *r.Name
}
