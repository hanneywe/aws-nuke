package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/amplify"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AmplifyWebhookResource = "AmplifyWebhook"

func init() {
	registry.Register(&registry.Registration{
		Name:     AmplifyWebhookResource,
		Scope:    nuke.Account,
		Resource: &AmplifyWebhook{},
		Lister:   &AmplifyWebhookLister{},
	})
}

type AmplifyWebhookLister struct {
	svc AmplifyClient
}

func (l *AmplifyWebhookLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = amplify.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	appPaginator := amplify.NewListAppsPaginator(svc, &amplify.ListAppsInput{})
	for appPaginator.HasMorePages() {
		appResp, err := appPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range appResp.Apps {
			params := &amplify.ListWebhooksInput{AppId: appResp.Apps[i].AppId}
			for {
				resp, err := svc.ListWebhooks(ctx, params)
				if err != nil {
					return nil, err
				}
				for _, wh := range resp.Webhooks {
					resources = append(resources, &AmplifyWebhook{
						svc:        svc,
						WebhookID:  wh.WebhookId,
						WebhookArn: wh.WebhookArn,
						AppID:      appResp.Apps[i].AppId,
					})
				}
				if resp.NextToken == nil {
					break
				}
				params.NextToken = resp.NextToken
			}
		}
	}
	return resources, nil
}

type AmplifyWebhook struct {
	svc        AmplifyClient
	WebhookID  *string `property:"name=WebhookId"`
	WebhookArn *string
	AppID      *string `property:"name=AppId"`
}

func (r *AmplifyWebhook) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteWebhook(ctx, &amplify.DeleteWebhookInput{
		WebhookId: r.WebhookID,
	})
	return err
}

func (r *AmplifyWebhook) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AmplifyWebhook) String() string {
	return *r.WebhookID
}
