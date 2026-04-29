package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmtypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SSMDefaultPatchBaselineResource = "SSMDefaultPatchBaseline"

func init() {
	registry.Register(&registry.Registration{
		Name:     SSMDefaultPatchBaselineResource,
		Scope:    nuke.Account,
		Resource: &SSMDefaultPatchBaseline{},
		Lister:   &SSMDefaultPatchBaselineLister{},
	})
}

type SSMDefaultPatchBaselineLister struct {
	svc SSMV2Client
}

func (l *SSMDefaultPatchBaselineLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = ssm.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// Find custom baselines that are set as default
	params := &ssm.DescribePatchBaselinesInput{
		Filters: []ssmtypes.PatchOrchestratorFilter{
			{
				Key:    aws.String("OWNER"),
				Values: []string{"Self"},
			},
		},
	}

	for {
		resp, err := svc.DescribePatchBaselines(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.BaselineIdentities {
			if resp.BaselineIdentities[i].DefaultBaseline {
				resources = append(resources, &SSMDefaultPatchBaseline{
					svc:             svc,
					BaselineID:      resp.BaselineIdentities[i].BaselineId,
					BaselineName:    resp.BaselineIdentities[i].BaselineName,
					OperatingSystem: resp.BaselineIdentities[i].OperatingSystem,
				})
			}
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type SSMDefaultPatchBaseline struct {
	svc             SSMV2Client
	BaselineID      *string
	BaselineName    *string
	OperatingSystem ssmtypes.OperatingSystem
}

func (r *SSMDefaultPatchBaseline) Remove(ctx context.Context) error {
	// Find the AWS-provided default baseline for this OS
	awsBaselineID, err := r.findAWSDefaultBaseline(ctx)
	if err != nil {
		return err
	}

	_, err = r.svc.RegisterDefaultPatchBaseline(ctx, &ssm.RegisterDefaultPatchBaselineInput{
		BaselineId: awsBaselineID,
	})
	return err
}

func (r *SSMDefaultPatchBaseline) findAWSDefaultBaseline(ctx context.Context) (*string, error) {
	params := &ssm.DescribePatchBaselinesInput{
		Filters: []ssmtypes.PatchOrchestratorFilter{
			{
				Key:    aws.String("OWNER"),
				Values: []string{"AWS"},
			},
			{
				Key:    aws.String("OPERATING_SYSTEM"),
				Values: []string{string(r.OperatingSystem)},
			},
		},
	}

	for {
		resp, err := r.svc.DescribePatchBaselines(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.BaselineIdentities {
			name := aws.ToString(resp.BaselineIdentities[i].BaselineName)
			if strings.HasPrefix(name, "AWS-") {
				return resp.BaselineIdentities[i].BaselineId, nil
			}
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return nil, fmt.Errorf("could not find AWS-provided default patch baseline for %s", r.OperatingSystem)
}

func (r *SSMDefaultPatchBaseline) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SSMDefaultPatchBaseline) String() string {
	return *r.BaselineID
}
