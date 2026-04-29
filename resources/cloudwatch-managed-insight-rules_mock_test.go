package resources

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	cloudwatchtypes "github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

func Test_Mock_CloudWatchManagedInsightRules_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchV2Client)

	mockClient.On("DescribeInsightRules", mock.Anything, mock.Anything).
		Return(&cloudwatch.DescribeInsightRulesOutput{
			InsightRules: []cloudwatchtypes.InsightRule{
				{
					Name:        ptr.String("managed-rule-1"),
					State:       ptr.String("ENABLED"),
					ManagedRule: ptr.Bool(true),
					Schema:      ptr.String(`{"Name": "ServiceLogRule", "Version": 1}`),
					Definition:  ptr.String("{}"),
				},
				{
					Name:        ptr.String("user-rule-1"),
					State:       ptr.String("ENABLED"),
					ManagedRule: ptr.Bool(false),
					Schema:      ptr.String(`{"Name": "CloudWatchLogRule", "Version": 1}`),
					Definition:  ptr.String("{}"),
				},
			},
		}, nil)

	lister := &CloudWatchManagedInsightRulesLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*CloudWatchManagedInsightRules)
	a.Equal("managed-rule-1", *r.Name)
	a.Equal("ENABLED", *r.State)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchManagedInsightRules_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchV2Client)

	mockClient.On("DescribeInsightRules", mock.Anything, mock.Anything).
		Return(&cloudwatch.DescribeInsightRulesOutput{
			InsightRules: []cloudwatchtypes.InsightRule{},
		}, nil)

	lister := &CloudWatchManagedInsightRulesLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchManagedInsightRules_List_NoManaged(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchV2Client)

	mockClient.On("DescribeInsightRules", mock.Anything, mock.Anything).
		Return(&cloudwatch.DescribeInsightRulesOutput{
			InsightRules: []cloudwatchtypes.InsightRule{
				{
					Name:        ptr.String("user-rule-1"),
					State:       ptr.String("ENABLED"),
					ManagedRule: ptr.Bool(false),
					Schema:      ptr.String(`{"Name": "CloudWatchLogRule", "Version": 1}`),
					Definition:  ptr.String("{}"),
				},
			},
		}, nil)

	lister := &CloudWatchManagedInsightRulesLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), &nuke.ListerOpts{})
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchManagedInsightRules_Filter_Disabled(t *testing.T) {
	a := assert.New(t)

	r := &CloudWatchManagedInsightRules{
		Name:  ptr.String("managed-rule-1"),
		State: ptr.String("DISABLED"),
	}

	err := r.Filter()
	a.Error(err)
	a.Contains(err.Error(), "already disabled")
}

func Test_Mock_CloudWatchManagedInsightRules_Filter_Enabled(t *testing.T) {
	a := assert.New(t)

	r := &CloudWatchManagedInsightRules{
		Name:  ptr.String("managed-rule-1"),
		State: ptr.String("ENABLED"),
	}

	err := r.Filter()
	a.Nil(err)
}

func Test_Mock_CloudWatchManagedInsightRules_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockCloudWatchV2Client)

	r := &CloudWatchManagedInsightRules{
		svc:  mockClient,
		Name: ptr.String("managed-rule-1"),
	}

	mockClient.On("DeleteInsightRules", mock.Anything, &cloudwatch.DeleteInsightRulesInput{
		RuleNames: []string{"managed-rule-1"},
	}).Return(&cloudwatch.DeleteInsightRulesOutput{}, nil)

	err := r.Remove(context.TODO())
	a.NoError(err)
	mockClient.AssertExpectations(t)
}

func Test_Mock_CloudWatchManagedInsightRules_Properties(t *testing.T) {
	a := assert.New(t)

	r := CloudWatchManagedInsightRules{
		Name:  ptr.String("managed-rule-1"),
		State: ptr.String("ENABLED"),
	}

	props := r.Properties()
	a.Equal("managed-rule-1", props.Get("Name"))
	a.Equal("ENABLED", props.Get("State"))
}

func Test_Mock_CloudWatchManagedInsightRules_String(t *testing.T) {
	a := assert.New(t)
	r := CloudWatchManagedInsightRules{
		Name: ptr.String("managed-rule-1"),
	}
	a.Equal("managed-rule-1", r.String())
}
