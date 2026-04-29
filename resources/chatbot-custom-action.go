package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/chatbot"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ChatbotCustomActionResource = "ChatbotCustomAction"

func init() {
	registry.Register(&registry.Registration{
		Name:     ChatbotCustomActionResource,
		Scope:    nuke.Account,
		Resource: &ChatbotCustomAction{},
		Lister:   &ChatbotCustomActionLister{},
	})
}

type ChatbotCustomActionLister struct {
	svc ChatbotClient
}

func (l *ChatbotCustomActionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = chatbot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := chatbot.NewListCustomActionsPaginator(svc, &chatbot.ListCustomActionsInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, arn := range resp.CustomActions {
			descResp, err := svc.GetCustomAction(ctx, &chatbot.GetCustomActionInput{
				CustomActionArn: &arn,
			})
			if err != nil {
				return nil, err
			}

			resources = append(resources, &ChatbotCustomAction{
				svc:             svc,
				CustomActionArn: descResp.CustomAction.CustomActionArn,
				ActionName:      descResp.CustomAction.ActionName,
			})
		}
	}

	return resources, nil
}

type ChatbotCustomAction struct {
	svc             ChatbotClient
	CustomActionArn *string
	ActionName      *string
}

func (r *ChatbotCustomAction) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteCustomAction(ctx, &chatbot.DeleteCustomActionInput{
		CustomActionArn: r.CustomActionArn,
	})
	return err
}

func (r *ChatbotCustomAction) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ChatbotCustomAction) String() string {
	return *r.CustomActionArn
}
