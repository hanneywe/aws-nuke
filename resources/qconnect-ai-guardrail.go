package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/qconnect"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const QConnectAIGuardrailResource = "QConnectAIGuardrail"

func init() {
	registry.Register(&registry.Registration{
		Name:     QConnectAIGuardrailResource,
		Scope:    nuke.Account,
		Resource: &QConnectAIGuardrail{},
		Lister:   &QConnectAIGuardrailLister{},
	})
}

type QConnectAIGuardrailLister struct {
	svc QConnectClient
}

func (l *QConnectAIGuardrailLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = qconnect.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	assistantPaginator := qconnect.NewListAssistantsPaginator(svc, &qconnect.ListAssistantsInput{})
	for assistantPaginator.HasMorePages() {
		assistantResp, err := assistantPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range assistantResp.AssistantSummaries {
			guardrailPaginator := qconnect.NewListAIGuardrailsPaginator(svc, &qconnect.ListAIGuardrailsInput{
				AssistantId: assistantResp.AssistantSummaries[i].AssistantId,
			})
			for guardrailPaginator.HasMorePages() {
				guardrailResp, err := guardrailPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for j := range guardrailResp.AiGuardrailSummaries {
					resources = append(resources, &QConnectAIGuardrail{
						svc:           svc,
						AssistantID:   assistantResp.AssistantSummaries[i].AssistantId,
						AIGuardrailID: guardrailResp.AiGuardrailSummaries[j].AiGuardrailId,
						Name:          guardrailResp.AiGuardrailSummaries[j].Name,
					})
				}
			}
		}
	}

	return resources, nil
}

type QConnectAIGuardrail struct {
	svc           QConnectClient
	AssistantID   *string `property:"name=AssistantId"`
	AIGuardrailID *string `property:"name=AIGuardrailId"`
	Name          *string
}

func (r *QConnectAIGuardrail) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAIGuardrail(ctx, &qconnect.DeleteAIGuardrailInput{
		AiGuardrailId: r.AIGuardrailID,
		AssistantId:   r.AssistantID,
	})
	return err
}

func (r *QConnectAIGuardrail) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *QConnectAIGuardrail) String() string {
	return *r.AIGuardrailID
}
