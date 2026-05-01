package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockagentruntime"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const BedrockAgentRuntimeSessionResource = "BedrockAgentRuntimeSession"

func init() {
	registry.Register(&registry.Registration{
		Name:     BedrockAgentRuntimeSessionResource,
		Scope:    nuke.Account,
		Resource: &BedrockAgentRuntimeSession{},
		Lister:   &BedrockAgentRuntimeSessionLister{},
	})
}

type BedrockAgentRuntimeSessionLister struct {
	svc BedrockagentruntimeClient
}

func (l *BedrockAgentRuntimeSessionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = bedrockagentruntime.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	paginator := bedrockagentruntime.NewListSessionsPaginator(svc, &bedrockagentruntime.ListSessionsInput{
		MaxResults: aws.Int32(100),
	})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range resp.SessionSummaries {
			session := &resp.SessionSummaries[i]
			status := string(session.SessionStatus)
			resources = append(resources, &BedrockAgentRuntimeSession{
				svc:           svc,
				SessionID:     session.SessionId,
				SessionArn:    session.SessionArn,
				SessionStatus: &status,
				CreatedAt:     session.CreatedAt,
				LastUpdatedAt: session.LastUpdatedAt,
			})
		}
	}

	return resources, nil
}

type BedrockAgentRuntimeSession struct {
	svc           BedrockagentruntimeClient
	SessionID     *string
	SessionArn    *string
	SessionStatus *string
	CreatedAt     *time.Time
	LastUpdatedAt *time.Time
}

func (r *BedrockAgentRuntimeSession) Remove(ctx context.Context) error {
	if r.SessionStatus != nil && *r.SessionStatus == "ACTIVE" {
		_, err := r.svc.EndSession(ctx, &bedrockagentruntime.EndSessionInput{
			SessionIdentifier: r.SessionID,
		})
		if err != nil {
			return err
		}
	}

	_, err := r.svc.DeleteSession(ctx, &bedrockagentruntime.DeleteSessionInput{
		SessionIdentifier: r.SessionID,
	})
	return err
}

func (r *BedrockAgentRuntimeSession) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *BedrockAgentRuntimeSession) String() string {
	return *r.SessionID
}
