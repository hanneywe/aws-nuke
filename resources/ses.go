package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// SESClient is the interface for the ses SDK client methods.
type SESClient interface {
	ListReceiptRuleSets(ctx context.Context, params *ses.ListReceiptRuleSetsInput,
		optFns ...func(*ses.Options)) (*ses.ListReceiptRuleSetsOutput, error)
	DescribeReceiptRuleSet(ctx context.Context, params *ses.DescribeReceiptRuleSetInput,
		optFns ...func(*ses.Options)) (*ses.DescribeReceiptRuleSetOutput, error)
	DeleteReceiptRule(ctx context.Context, params *ses.DeleteReceiptRuleInput,
		optFns ...func(*ses.Options)) (*ses.DeleteReceiptRuleOutput, error)
}
