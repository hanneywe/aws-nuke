package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lexmodelsv2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const LexV2BotResource = "LexV2Bot"

func init() {
	registry.Register(&registry.Registration{
		Name:     LexV2BotResource,
		Scope:    nuke.Account,
		Resource: &LexV2Bot{},
		Lister:   &LexV2BotLister{},
	})
}

type LexV2BotLister struct {
	svc LexV2Client
}

func (l *LexV2BotLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = lexmodelsv2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &lexmodelsv2.ListBotsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		output, err := svc.ListBots(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, botSummary := range output.BotSummaries {
			resources = append(resources, &LexV2Bot{
				svc:     svc,
				BotID:   botSummary.BotId,
				BotName: botSummary.BotName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

type LexV2Bot struct {
	svc     LexV2Client
	BotID   *string `property:"name=BotId"`
	BotName *string
}

func (r *LexV2Bot) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteBot(ctx, &lexmodelsv2.DeleteBotInput{
		BotId:                  r.BotID,
		SkipResourceInUseCheck: true,
	})
	return err
}

func (r *LexV2Bot) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *LexV2Bot) String() string {
	return *r.BotName
}
