package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/qconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const QConnectKnowledgeBaseResource = "QConnectKnowledgeBase"

func init() {
	registry.Register(&registry.Registration{
		Name:     QConnectKnowledgeBaseResource,
		Scope:    nuke.Account,
		Resource: &QConnectKnowledgeBase{},
		Lister:   &QConnectKnowledgeBaseLister{},
	})
}

type QConnectKnowledgeBaseLister struct {
	svc QConnectClient
}

func (l *QConnectKnowledgeBaseLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = qconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	paginator := qconnect.NewListKnowledgeBasesPaginator(svc, &qconnect.ListKnowledgeBasesInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range resp.KnowledgeBaseSummaries {
			resources = append(resources, &QConnectKnowledgeBase{
				svc:             svc,
				KnowledgeBaseID: resp.KnowledgeBaseSummaries[i].KnowledgeBaseId,
				Name:            resp.KnowledgeBaseSummaries[i].Name,
			})
		}
	}

	return resources, nil
}

type QConnectKnowledgeBase struct {
	svc             QConnectClient
	KnowledgeBaseID *string `property:"name=KnowledgeBaseId"`
	Name            *string
}

func (r *QConnectKnowledgeBase) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteKnowledgeBase(ctx, &qconnect.DeleteKnowledgeBaseInput{
		KnowledgeBaseId: r.KnowledgeBaseID,
	})
	return err
}

func (r *QConnectKnowledgeBase) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *QConnectKnowledgeBase) String() string {
	return *r.KnowledgeBaseID
}
