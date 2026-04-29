package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSV2Client interface {
	ListTopics(ctx context.Context, params *sns.ListTopicsInput,
		optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error)
	GetDataProtectionPolicy(ctx context.Context, params *sns.GetDataProtectionPolicyInput,
		optFns ...func(*sns.Options)) (*sns.GetDataProtectionPolicyOutput, error)
	PutDataProtectionPolicy(ctx context.Context, params *sns.PutDataProtectionPolicyInput,
		optFns ...func(*sns.Options)) (*sns.PutDataProtectionPolicyOutput, error)
}
