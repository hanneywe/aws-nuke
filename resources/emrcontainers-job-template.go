package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/emrcontainers"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EMRContainersJobTemplateResource = "EMRContainersJobTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     EMRContainersJobTemplateResource,
		Scope:    nuke.Account,
		Resource: &EMRContainersJobTemplate{},
		Lister:   &EMRContainersJobTemplateLister{},
	})
}

type EMRContainersJobTemplateLister struct {
	svc EMRContainersClient
}

func (l *EMRContainersJobTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = emrcontainers.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	var nextToken *string

	for {
		resp, err := svc.ListJobTemplates(ctx, &emrcontainers.ListJobTemplatesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, t := range resp.Templates {
			resources = append(resources, &EMRContainersJobTemplate{
				svc:  svc,
				ID:   t.Id,
				Name: t.Name,
				ARN:  t.Arn,
				Tags: t.Tags,
			})
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return resources, nil
}

type EMRContainersJobTemplate struct {
	svc  EMRContainersClient
	ID   *string
	Name *string
	ARN  *string
	Tags map[string]string
}

func (r *EMRContainersJobTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteJobTemplate(ctx, &emrcontainers.DeleteJobTemplateInput{
		Id: r.ID,
	})
	return err
}

func (r *EMRContainersJobTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EMRContainersJobTemplate) String() string {
	return *r.Name
}
