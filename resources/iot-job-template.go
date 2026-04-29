package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTJobTemplateResource = "IoTJobTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTJobTemplateResource,
		Scope:    nuke.Account,
		Resource: &IoTJobTemplate{},
		Lister:   &IoTJobTemplateLister{},
	})
}

type IoTJobTemplateLister struct {
	svc IoTClient
}

func (l *IoTJobTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := iot.NewListJobTemplatesPaginator(svc, &iot.ListJobTemplatesInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, jt := range resp.JobTemplates {
			resources = append(resources, &IoTJobTemplate{
				svc:            svc,
				JobTemplateID:  jt.JobTemplateId,
				JobTemplateArn: jt.JobTemplateArn,
			})
		}
	}

	return resources, nil
}

type IoTJobTemplate struct {
	svc            IoTClient
	JobTemplateID  *string `property:"name=JobTemplateId"`
	JobTemplateArn *string
}

func (r *IoTJobTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteJobTemplate(ctx, &iot.DeleteJobTemplateInput{
		JobTemplateId: r.JobTemplateID,
	})
	return err
}

func (r *IoTJobTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTJobTemplate) String() string {
	return *r.JobTemplateID
}
