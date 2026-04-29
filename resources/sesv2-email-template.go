package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SESv2EmailTemplateResource = "SESv2EmailTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     SESv2EmailTemplateResource,
		Scope:    nuke.Account,
		Resource: &SESv2EmailTemplate{},
		Lister:   &SESv2EmailTemplateLister{},
	})
}

type SESv2EmailTemplateLister struct {
	svc SESv2Client
}

func (l *SESv2EmailTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = sesv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := sesv2.NewListEmailTemplatesPaginator(svc, &sesv2.ListEmailTemplatesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, tmpl := range resp.TemplatesMetadata {
			resources = append(resources, &SESv2EmailTemplate{
				svc:          svc,
				TemplateName: tmpl.TemplateName,
			})
		}
	}
	return resources, nil
}

type SESv2EmailTemplate struct {
	svc          SESv2Client
	TemplateName *string
}

func (r *SESv2EmailTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEmailTemplate(ctx, &sesv2.DeleteEmailTemplateInput{
		TemplateName: r.TemplateName,
	})
	return err
}

func (r *SESv2EmailTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SESv2EmailTemplate) String() string {
	return *r.TemplateName
}
