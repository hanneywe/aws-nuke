package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ses"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockSESClient struct {
	mock.Mock
}

func (m *mockSESClient) ListReceiptRuleSets(
	ctx context.Context, params *ses.ListReceiptRuleSetsInput,
	_ ...func(*ses.Options),
) (*ses.ListReceiptRuleSetsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ses.ListReceiptRuleSetsOutput), args.Error(1)
}

func (m *mockSESClient) DescribeReceiptRuleSet(
	ctx context.Context, params *ses.DescribeReceiptRuleSetInput,
	_ ...func(*ses.Options),
) (*ses.DescribeReceiptRuleSetOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ses.DescribeReceiptRuleSetOutput), args.Error(1)
}

func (m *mockSESClient) DeleteReceiptRule(
	ctx context.Context, params *ses.DeleteReceiptRuleInput,
	_ ...func(*ses.Options),
) (*ses.DeleteReceiptRuleOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ses.DeleteReceiptRuleOutput), args.Error(1)
}

var testSesListerOpts = &nuke.ListerOpts{}
