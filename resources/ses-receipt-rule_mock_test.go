package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	sestypes "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

func Test_Mock_SESReceiptRule_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESClient)

	mockClient.On("ListReceiptRuleSets", mock.Anything, mock.Anything).
		Return(&ses.ListReceiptRuleSetsOutput{
			RuleSets: []sestypes.ReceiptRuleSetMetadata{
				{Name: ptr.String("test-rulesetname")},
			},
		}, nil)

	mockClient.On("DescribeReceiptRuleSet", mock.Anything, mock.Anything).
		Return(&ses.DescribeReceiptRuleSetOutput{
			Rules: []sestypes.ReceiptRule{
				{Name: ptr.String("test-rulename"), Enabled: true},
			},
		}, nil)

	lister := &SESReceiptRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSesListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*SESReceiptRule)
	a.Equal("test-rulesetname", *r.RuleSetName)
	a.Equal("test-rulename", *r.RuleName)
	a.Equal(true, r.Enabled)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESReceiptRule_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESClient)

	mockClient.On("ListReceiptRuleSets", mock.Anything, mock.Anything).
		Return(&ses.ListReceiptRuleSetsOutput{
			RuleSets: []sestypes.ReceiptRuleSetMetadata{},
		}, nil)

	lister := &SESReceiptRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testSesListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESReceiptRule_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockSESClient)

	r := &SESReceiptRule{
		svc:         mockClient,
		RuleName:    ptr.String("test-rulename"),
		RuleSetName: ptr.String("test-rulesetname"),
	}

	mockClient.On("DeleteReceiptRule", mock.Anything,
		&ses.DeleteReceiptRuleInput{
			RuleName:    r.RuleName,
			RuleSetName: r.RuleSetName,
		}).Return(&ses.DeleteReceiptRuleOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_SESReceiptRule_Properties(t *testing.T) {
	a := assert.New(t)
	r := &SESReceiptRule{
		RuleSetName: ptr.String("test-rulesetname"),
		RuleName:    ptr.String("test-rulename"),
	}
	props := r.Properties()
	a.Equal("test-rulesetname", props.Get("RuleSetName"))
	a.Equal("test-rulename", props.Get("RuleName"))
}

func Test_Mock_SESReceiptRule_String(t *testing.T) {
	a := assert.New(t)
	r := &SESReceiptRule{
		RuleName: ptr.String("test-rulename"),
	}
	a.Equal("test-rulename", r.String())
}
