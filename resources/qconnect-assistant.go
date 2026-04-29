package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/qconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const QConnectAssistantResource = "QConnectAssistant"

func init() {
	registry.Register(&registry.Registration{
		Name:     QConnectAssistantResource,
		Scope:    nuke.Account,
		Resource: &QConnectAssistant{},
		Lister:   &QConnectAssistantLister{},
	})
}

type QConnectAssistantLister struct {
	svc QConnectClient
}

func (l *QConnectAssistantLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = qconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := qconnect.NewListAssistantsPaginator(svc, &qconnect.ListAssistantsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range resp.AssistantSummaries {
			resources = append(resources, &QConnectAssistant{
				svc:         svc,
				AssistantID: resp.AssistantSummaries[i].AssistantId,
				Name:        resp.AssistantSummaries[i].Name,
			})
		}
	}
	return resources, nil
}

type QConnectAssistant struct {
	svc         QConnectClient
	AssistantID *string `property:"name=AssistantId"`
	Name        *string
}

func (r *QConnectAssistant) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAssistant(ctx, &qconnect.DeleteAssistantInput{
		AssistantId: r.AssistantID,
	})
	return err
}

func (r *QConnectAssistant) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *QConnectAssistant) String() string {
	return *r.Name
}
