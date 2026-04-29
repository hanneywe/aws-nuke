package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CloudFormationGeneratedTemplateResource = "CloudFormationGeneratedTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     CloudFormationGeneratedTemplateResource,
		Scope:    nuke.Account,
		Resource: &CloudFormationGeneratedTemplate{},
		Lister:   &CloudFormationGeneratedTemplateLister{},
	})
}

type CloudFormationGeneratedTemplateLister struct {
	svc CloudFormationClient
}

func (l *CloudFormationGeneratedTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cloudformation.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &cloudformation.ListGeneratedTemplatesInput{}
	for {
		resp, err := svc.ListGeneratedTemplates(ctx, params)
		if err != nil {
			return nil, err
		}
		for _, t := range resp.Summaries {
			resources = append(resources, &CloudFormationGeneratedTemplate{
				svc:                   svc,
				GeneratedTemplateID:   t.GeneratedTemplateId,
				GeneratedTemplateName: t.GeneratedTemplateName,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type CloudFormationGeneratedTemplate struct {
	svc                   CloudFormationClient
	GeneratedTemplateID   *string `property:"name=GeneratedTemplateId"`
	GeneratedTemplateName *string
}

func (r *CloudFormationGeneratedTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteGeneratedTemplate(ctx, &cloudformation.DeleteGeneratedTemplateInput{
		GeneratedTemplateName: r.GeneratedTemplateName,
	})
	return err
}

func (r *CloudFormationGeneratedTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CloudFormationGeneratedTemplate) String() string {
	return *r.GeneratedTemplateName
}
