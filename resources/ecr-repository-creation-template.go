package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ECRRepositoryCreationTemplateResource = "ECRRepositoryCreationTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     ECRRepositoryCreationTemplateResource,
		Scope:    nuke.Account,
		Resource: &ECRRepositoryCreationTemplate{},
		Lister:   &ECRRepositoryCreationTemplateLister{},
	})
}

type ECRRepositoryCreationTemplateLister struct {
	svc ECRv2Client
}

func (l *ECRRepositoryCreationTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = ecr.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &ecr.DescribeRepositoryCreationTemplatesInput{}
	for {
		resp, err := svc.DescribeRepositoryCreationTemplates(ctx, params)
		if err != nil {
			return nil, err
		}
		for i := range resp.RepositoryCreationTemplates {
			resources = append(resources, &ECRRepositoryCreationTemplate{
				svc:    svc,
				Prefix: resp.RepositoryCreationTemplates[i].Prefix,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}
	return resources, nil
}

type ECRRepositoryCreationTemplate struct {
	svc    ECRv2Client
	Prefix *string
}

func (r *ECRRepositoryCreationTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRepositoryCreationTemplate(ctx, &ecr.DeleteRepositoryCreationTemplateInput{
		Prefix: r.Prefix,
	})
	return err
}

func (r *ECRRepositoryCreationTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ECRRepositoryCreationTemplate) String() string {
	return *r.Prefix
}
