package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/batch"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BatchJobDefinitionResource = "BatchJobDefinition"

func init() {
	registry.Register(&registry.Registration{
		Name:     BatchJobDefinitionResource,
		Scope:    nuke.Account,
		Resource: &BatchJobDefinition{},
		Lister:   &BatchJobDefinitionLister{},
	})
}

type BatchJobDefinitionLister struct {
	svc BatchClient
}

func (l *BatchJobDefinitionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = batch.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := batch.NewDescribeJobDefinitionsPaginator(svc, &batch.DescribeJobDefinitionsInput{
		Status: aws.String("ACTIVE"),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.JobDefinitions {
			jd := &resp.JobDefinitions[i]
			resources = append(resources, &BatchJobDefinition{
				svc:               svc,
				JobDefinitionArn:  jd.JobDefinitionArn,
				JobDefinitionName: jd.JobDefinitionName,
				Status:            jd.Status,
			})
		}
	}

	return resources, nil
}

type BatchJobDefinition struct {
	svc               BatchClient
	JobDefinitionArn  *string
	JobDefinitionName *string
	Status            *string
}

func (r *BatchJobDefinition) Filter() error {
	if r.Status != nil && *r.Status == "INACTIVE" {
		return fmt.Errorf("job definition is inactive")
	}
	return nil
}

func (r *BatchJobDefinition) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterJobDefinition(ctx, &batch.DeregisterJobDefinitionInput{
		JobDefinition: r.JobDefinitionArn,
	})
	return err
}

func (r *BatchJobDefinition) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BatchJobDefinition) String() string {
	return *r.JobDefinitionArn
}
