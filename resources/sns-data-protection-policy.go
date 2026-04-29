package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const SNSDataProtectionPolicyResource = "SNSDataProtectionPolicy"

func init() {
	registry.Register(&registry.Registration{
		Name:     SNSDataProtectionPolicyResource,
		Scope:    nuke.Account,
		Resource: &SNSDataProtectionPolicy{},
		Lister:   &SNSDataProtectionPolicyLister{},
	})
}

type SNSDataProtectionPolicyLister struct {
	svc SNSV2Client
}

func (l *SNSDataProtectionPolicyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = sns.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &sns.ListTopicsInput{}
	for {
		resp, err := svc.ListTopics(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.Topics {
			topicArn := resp.Topics[i].TopicArn
			policyResp, err := svc.GetDataProtectionPolicy(ctx, &sns.GetDataProtectionPolicyInput{
				ResourceArn: topicArn,
			})
			if err != nil {
				continue
			}

			if policyResp.DataProtectionPolicy == nil || *policyResp.DataProtectionPolicy == "" {
				continue
			}

			resources = append(resources, &SNSDataProtectionPolicy{
				svc:      svc,
				TopicArn: topicArn,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type SNSDataProtectionPolicy struct {
	svc      SNSV2Client
	TopicArn *string
}

func (r *SNSDataProtectionPolicy) Remove(ctx context.Context) error {
	_, err := r.svc.PutDataProtectionPolicy(ctx, &sns.PutDataProtectionPolicyInput{
		ResourceArn:          r.TopicArn,
		DataProtectionPolicy: aws.String(""),
	})
	return err
}

func (r *SNSDataProtectionPolicy) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *SNSDataProtectionPolicy) String() string {
	return *r.TopicArn
}
