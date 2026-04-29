package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configtypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
)

func Test_Mock_ConfigServiceOrganizationConfigRule_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeOrganizationConfigRules", mock.Anything, mock.Anything).
		Return(&configservice.DescribeOrganizationConfigRulesOutput{
			OrganizationConfigRules: []configtypes.OrganizationConfigRule{
				{
					OrganizationConfigRuleName: ptr.String("my-org-rule"),
					OrganizationConfigRuleArn:  ptr.String("arn:aws:config:us-east-1:123456789012:organization-config-rule/my-org-rule"),
				},
			},
		}, nil)
	lister := &ConfigServiceOrganizationConfigRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-org-rule", resources[0].(*ConfigServiceOrganizationConfigRule).String())
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceOrganizationConfigRule_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	mockClient.On("DescribeOrganizationConfigRules", mock.Anything, mock.Anything).
		Return(&configservice.DescribeOrganizationConfigRulesOutput{
			OrganizationConfigRules: []configtypes.OrganizationConfigRule{},
		}, nil)
	lister := &ConfigServiceOrganizationConfigRuleLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testConfigServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceOrganizationConfigRule_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockConfigServiceClient)
	r := &ConfigServiceOrganizationConfigRule{
		svc:  mockClient,
		Name: ptr.String("my-org-rule"),
	}
	mockClient.On("DeleteOrganizationConfigRule", mock.Anything, &configservice.DeleteOrganizationConfigRuleInput{
		OrganizationConfigRuleName: r.Name,
	}).Return(&configservice.DeleteOrganizationConfigRuleOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ConfigServiceOrganizationConfigRule_Properties(t *testing.T) {
	a := assert.New(t)
	r := ConfigServiceOrganizationConfigRule{
		Name: ptr.String("my-org-rule"),
	}
	a.Equal("my-org-rule", r.Properties().Get("Name"))
}

func Test_Mock_ConfigServiceOrganizationConfigRule_String(t *testing.T) {
	a := assert.New(t)
	r := &ConfigServiceOrganizationConfigRule{
		Name: ptr.String("my-org-rule"),
	}
	a.Equal("my-org-rule", r.String())
}
