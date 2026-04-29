package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/glue"
	gluetypes "github.com/aws/aws-sdk-go-v2/service/glue/types"
)

func Test_Mock_GlueDataQualityRuleset_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListDataQualityRulesets", mock.Anything, mock.Anything).
		Return(&glue.ListDataQualityRulesetsOutput{
			Rulesets: []gluetypes.DataQualityRulesetListDetails{
				{Name: ptr.String("my-ruleset"), Description: ptr.String("test desc")},
			},
		}, nil)
	lister := &GlueDataQualityRulesetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-ruleset", resources[0].(*GlueDataQualityRuleset).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataQualityRuleset_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	mockClient.On("ListDataQualityRulesets", mock.Anything, mock.Anything).
		Return(&glue.ListDataQualityRulesetsOutput{Rulesets: []gluetypes.DataQualityRulesetListDetails{}}, nil)
	lister := &GlueDataQualityRulesetLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testGlueV2ListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataQualityRuleset_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockGlueV2Client)
	r := &GlueDataQualityRuleset{svc: mockClient, Name: ptr.String("my-ruleset")}
	mockClient.On("DeleteDataQualityRuleset", mock.Anything, &glue.DeleteDataQualityRulesetInput{
		Name: r.Name,
	}).Return(&glue.DeleteDataQualityRulesetOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_GlueDataQualityRuleset_Properties(t *testing.T) {
	a := assert.New(t)
	r := GlueDataQualityRuleset{Name: ptr.String("my-ruleset"), Description: ptr.String("test desc")}
	a.Equal("my-ruleset", r.Properties().Get("Name"))
	a.Equal("test desc", r.Properties().Get("Description"))
}

func Test_Mock_GlueDataQualityRuleset_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-ruleset", (&GlueDataQualityRuleset{Name: ptr.String("my-ruleset")}).String())
}
