package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/kendraranking"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const KendraRankingRescoreExecutionPlanResource = "KendraRankingRescoreExecutionPlan"

func init() {
	registry.Register(&registry.Registration{
		Name:     KendraRankingRescoreExecutionPlanResource,
		Scope:    nuke.Account,
		Resource: &KendraRankingRescoreExecutionPlan{},
		Lister:   &KendraRankingRescoreExecutionPlanLister{},
	})
}

type KendraRankingRescoreExecutionPlanLister struct {
	svc KendraRankingClient
}

func (l *KendraRankingRescoreExecutionPlanLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = kendraranking.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := kendraranking.NewListRescoreExecutionPlansPaginator(svc, &kendraranking.ListRescoreExecutionPlansInput{})

	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, plan := range resp.SummaryItems {
			resources = append(resources, &KendraRankingRescoreExecutionPlan{
				svc:  svc,
				ID:   plan.Id,
				Name: plan.Name,
			})
		}
	}

	return resources, nil
}

type KendraRankingRescoreExecutionPlan struct {
	svc  KendraRankingClient
	ID   *string `property:"name=Id"`
	Name *string
}

func (r *KendraRankingRescoreExecutionPlan) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRescoreExecutionPlan(ctx, &kendraranking.DeleteRescoreExecutionPlanInput{
		Id: r.ID,
	})
	return err
}

func (r *KendraRankingRescoreExecutionPlan) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *KendraRankingRescoreExecutionPlan) String() string {
	return *r.Name
}
