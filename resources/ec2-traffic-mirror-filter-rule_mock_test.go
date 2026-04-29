package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func Test_Mock_EC2TrafficMirrorFilterRule_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorFilters", mock.Anything, mock.Anything).
		Return(&ec2.DescribeTrafficMirrorFiltersOutput{
			TrafficMirrorFilters: []ec2types.TrafficMirrorFilter{
				{
					TrafficMirrorFilterId: ptr.String("tmf-11111111111111111"),
				},
			},
		}, nil)

	mockClient.
		On("DescribeTrafficMirrorFilterRules", mock.Anything, mock.Anything).
		Return(&ec2.DescribeTrafficMirrorFilterRulesOutput{
			TrafficMirrorFilterRules: []ec2types.TrafficMirrorFilterRule{
				{
					TrafficMirrorFilterRuleId: ptr.String("tmfr-11111111111111111"),
					TrafficMirrorFilterId:     ptr.String("tmf-11111111111111111"),
					TrafficDirection:          ec2types.TrafficDirectionIngress,
					RuleNumber:                ptr.Int32(100),
					RuleAction:                ec2types.TrafficMirrorRuleActionAccept,
					Tags: []ec2types.Tag{
						{Key: ptr.String("Name"), Value: ptr.String("test-rule")},
					},
				},
			},
		}, nil)

	lister := &EC2TrafficMirrorFilterRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	rule := resources[0].(*EC2TrafficMirrorFilterRule)
	a.Equal("tmfr-11111111111111111", *rule.TrafficMirrorFilterRuleID)
	a.Equal("tmf-11111111111111111", *rule.TrafficMirrorFilterID)
	a.Equal("ingress", rule.TrafficDirection)
	a.Equal(int32(100), *rule.RuleNumber)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilterRule_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorFilters", mock.Anything, mock.Anything).
		Return(&ec2.DescribeTrafficMirrorFiltersOutput{
			TrafficMirrorFilters: []ec2types.TrafficMirrorFilter{},
		}, nil)

	lister := &EC2TrafficMirrorFilterRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilterRule_List_MultipleFilters(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	mockClient.
		On("DescribeTrafficMirrorFilters", mock.Anything, mock.Anything).
		Return(&ec2.DescribeTrafficMirrorFiltersOutput{
			TrafficMirrorFilters: []ec2types.TrafficMirrorFilter{
				{TrafficMirrorFilterId: ptr.String("tmf-1")},
				{TrafficMirrorFilterId: ptr.String("tmf-2")},
			},
		}, nil)

	mockClient.
		On("DescribeTrafficMirrorFilterRules", mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorFilterRulesInput) bool {
				return *input.TrafficMirrorFilterId == "tmf-1"
			})).
		Return(&ec2.DescribeTrafficMirrorFilterRulesOutput{
			TrafficMirrorFilterRules: []ec2types.TrafficMirrorFilterRule{
				{
					TrafficMirrorFilterRuleId: ptr.String("tmfr-1"),
					TrafficMirrorFilterId:     ptr.String("tmf-1"),
					TrafficDirection:          ec2types.TrafficDirectionIngress,
					RuleNumber:                ptr.Int32(100),
					RuleAction:                ec2types.TrafficMirrorRuleActionAccept,
				},
			},
		}, nil)

	mockClient.
		On("DescribeTrafficMirrorFilterRules", mock.Anything,
			mock.MatchedBy(func(input *ec2.DescribeTrafficMirrorFilterRulesInput) bool {
				return *input.TrafficMirrorFilterId == "tmf-2"
			})).
		Return(&ec2.DescribeTrafficMirrorFilterRulesOutput{
			TrafficMirrorFilterRules: []ec2types.TrafficMirrorFilterRule{
				{
					TrafficMirrorFilterRuleId: ptr.String("tmfr-2"),
					TrafficMirrorFilterId:     ptr.String("tmf-2"),
					TrafficDirection:          ec2types.TrafficDirectionEgress,
					RuleNumber:                ptr.Int32(200),
					RuleAction:                ec2types.TrafficMirrorRuleActionReject,
				},
				{
					TrafficMirrorFilterRuleId: ptr.String("tmfr-3"),
					TrafficMirrorFilterId:     ptr.String("tmf-2"),
					TrafficDirection:          ec2types.TrafficDirectionIngress,
					RuleNumber:                ptr.Int32(300),
					RuleAction:                ec2types.TrafficMirrorRuleActionAccept,
				},
			},
		}, nil)

	lister := &EC2TrafficMirrorFilterRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testEC2ClientListerOpts)
	a.NoError(err)
	a.Len(resources, 3)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilterRule_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockEC2Client)

	rule := &EC2TrafficMirrorFilterRule{
		svc:                       mockClient,
		TrafficMirrorFilterRuleID: ptr.String("tmfr-11111111111111111"),
	}

	mockClient.
		On("DeleteTrafficMirrorFilterRule", mock.Anything, &ec2.DeleteTrafficMirrorFilterRuleInput{
			TrafficMirrorFilterRuleId: rule.TrafficMirrorFilterRuleID,
		}).
		Return(&ec2.DeleteTrafficMirrorFilterRuleOutput{}, nil)

	err := rule.Remove(context.TODO())
	a.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_EC2TrafficMirrorFilterRule_Properties(t *testing.T) {
	a := assert.New(t)

	rule := EC2TrafficMirrorFilterRule{
		TrafficMirrorFilterRuleID: ptr.String("tmfr-11111111111111111"),
		TrafficMirrorFilterID:     ptr.String("tmf-11111111111111111"),
		TrafficDirection:          "ingress",
		RuleNumber:                ptr.Int32(100),
		RuleAction:                "accept",
		Tags: []ec2types.Tag{
			{Key: ptr.String("Env"), Value: ptr.String("prod")},
		},
	}

	props := rule.Properties()
	a.Equal("tmfr-11111111111111111", props.Get("TrafficMirrorFilterRuleId"))
	a.Equal("tmf-11111111111111111", props.Get("TrafficMirrorFilterId"))
	a.Equal("ingress", props.Get("TrafficDirection"))
	a.Equal("accept", props.Get("RuleAction"))
	a.Equal("prod", props.Get("tag:Env"))
}

func Test_Mock_EC2TrafficMirrorFilterRule_String(t *testing.T) {
	a := assert.New(t)

	rule := EC2TrafficMirrorFilterRule{
		TrafficMirrorFilterRuleID: ptr.String("tmfr-11111111111111111"),
		TrafficMirrorFilterID:     ptr.String("tmf-11111111111111111"),
	}

	a.Equal("tmf-11111111111111111 -> tmfr-11111111111111111", rule.String())
}
