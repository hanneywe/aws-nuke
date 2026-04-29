package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/kendraranking"
	kendrarankingtypes "github.com/aws/aws-sdk-go-v2/service/kendraranking/types"
)

func Test_Mock_KendraRankingRescoreExecutionPlan_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKendraRankingClient)

	mockClient.On("ListRescoreExecutionPlans", mock.Anything, mock.Anything).
		Return(&kendraranking.ListRescoreExecutionPlansOutput{
			SummaryItems: []kendrarankingtypes.RescoreExecutionPlanSummary{
				{
					Id:   ptr.String("plan-12345"),
					Name: ptr.String("my-plan"),
				},
			},
		}, nil)

	lister := &KendraRankingRescoreExecutionPlanLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKendraRankingListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	plan := resources[0].(*KendraRankingRescoreExecutionPlan)
	a.Equal("plan-12345", *plan.ID)
	a.Equal("my-plan", *plan.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KendraRankingRescoreExecutionPlan_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKendraRankingClient)

	mockClient.On("ListRescoreExecutionPlans", mock.Anything, mock.Anything).
		Return(&kendraranking.ListRescoreExecutionPlansOutput{
			SummaryItems: []kendrarankingtypes.RescoreExecutionPlanSummary{},
		}, nil)

	lister := &KendraRankingRescoreExecutionPlanLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testKendraRankingListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_KendraRankingRescoreExecutionPlan_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockKendraRankingClient)

	plan := &KendraRankingRescoreExecutionPlan{
		svc: mockClient,
		ID:  ptr.String("plan-12345"),
	}

	mockClient.On("DeleteRescoreExecutionPlan", mock.Anything, &kendraranking.DeleteRescoreExecutionPlanInput{
		Id: plan.ID,
	}).Return(&kendraranking.DeleteRescoreExecutionPlanOutput{}, nil)

	a.NoError(plan.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_KendraRankingRescoreExecutionPlan_Properties(t *testing.T) {
	a := assert.New(t)

	plan := KendraRankingRescoreExecutionPlan{
		ID:   ptr.String("plan-12345"),
		Name: ptr.String("my-plan"),
	}

	props := plan.Properties()
	a.Equal("plan-12345", props.Get("Id"))
	a.Equal("my-plan", props.Get("Name"))
}

func Test_Mock_KendraRankingRescoreExecutionPlan_String(t *testing.T) {
	a := assert.New(t)
	plan := KendraRankingRescoreExecutionPlan{Name: ptr.String("my-plan")}
	a.Equal("my-plan", plan.String())
}
