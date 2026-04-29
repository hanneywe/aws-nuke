package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SSMOpsItemResource = "SSMOpsItem"

func init() {
	registry.Register(&registry.Registration{
		Name:     SSMOpsItemResource,
		Scope:    nuke.Account,
		Resource: &SSMOpsItem{},
		Lister:   &SSMOpsItemLister{},
	})
}

type SSMOpsItemLister struct {
	svc SSMV2Client
}

func (l *SSMOpsItemLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ssm.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &ssm.DescribeOpsItemsInput{
		OpsItemFilters: []ssmtypes.OpsItemFilter{
			{
				Key:      ssmtypes.OpsItemFilterKeyStatus,
				Values:   []string{"Open", "InProgress"},
				Operator: ssmtypes.OpsItemFilterOperatorEqual,
			},
		},
	}

	for {
		output, err := svc.DescribeOpsItems(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range output.OpsItemSummaries {
			opsItemSummary := &output.OpsItemSummaries[i]
			resources = append(resources, &SSMOpsItem{
				svc:       svc,
				OpsItemID: opsItemSummary.OpsItemId,
				Title:     opsItemSummary.Title,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type SSMOpsItem struct {
	svc       SSMV2Client
	OpsItemID *string `property:"name=OpsItemId"`
	Title     *string
}

func (r *SSMOpsItem) Remove(ctx context.Context) error {
	_, err := r.svc.UpdateOpsItem(ctx, &ssm.UpdateOpsItemInput{
		OpsItemId: r.OpsItemID,
		Status:    ssmtypes.OpsItemStatusResolved,
	})
	return err
}

func (r *SSMOpsItem) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SSMOpsItem) String() string {
	return *r.OpsItemID
}
