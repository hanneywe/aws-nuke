package resources

import (
	"context"
	"fmt"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	lattice "github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
)

func Test_Mock_VPCLatticeListenerRule_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{Id: ptr.String("svc-1"), Name: ptr.String("service-one")},
				},
			}, nil,
		)

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{
					{Id: ptr.String("l-1"), Name: ptr.String("listener-1")},
				},
			}, nil,
		)

	mockClient.
		On("ListRules", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListRulesOutput{
				Items: []lattice.RuleSummary{
					{
						Id:        ptr.String("r-1"),
						Arn:       ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:rule/r-1"),
						Name:      ptr.String("rule-1"),
						Priority:  ptr.Int32(10),
						IsDefault: ptr.Bool(false),
					},
				},
			}, nil,
		)

	lister := &VPCLatticeListenerRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	rule := resources[0].(*VPCLatticeListenerRule)
	assertions.Equal("rule-1", *rule.Name)
	assertions.Equal("svc-1", *rule.ServiceID)
	assertions.Equal("l-1", *rule.ListenerID)
	assertions.False(rule.IsDefault)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListenerRule_List_Empty_NoServices(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{},
			}, nil,
		)

	lister := &VPCLatticeListenerRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListenerRule_List_Empty_NoListeners(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{Id: ptr.String("svc-1"), Name: ptr.String("service-one")},
				},
			}, nil,
		)

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{},
			}, nil,
		)

	lister := &VPCLatticeListenerRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListenerRule_List_Empty_NoRules(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	mockClient.
		On("ListServices", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{Id: ptr.String("svc-1"), Name: ptr.String("service-one")},
				},
			}, nil,
		)

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{
					{Id: ptr.String("l-1"), Name: ptr.String("listener-1")},
				},
			}, nil,
		)

	mockClient.
		On("ListRules", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListRulesOutput{
				Items: []lattice.RuleSummary{},
			}, nil,
		)

	lister := &VPCLatticeListenerRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListenerRule_List_MultiPage_Services(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	// 101 services across 2 pages, each with 1 listener and 1 rule = 101 rules total
	firstPageServices := make([]lattice.ServiceSummary, 100)
	for i := range firstPageServices {
		firstPageServices[i] = lattice.ServiceSummary{
			Id:   ptr.String(fmt.Sprintf("svc-%d", i)),
			Name: ptr.String(fmt.Sprintf("service-%d", i)),
		}
	}

	mockClient.
		On(
			"ListServices",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServicesInput) bool {
				return input.NextToken == nil
			}),
		).
		Return(
			&vpclattice.ListServicesOutput{
				Items:     firstPageServices,
				NextToken: ptr.String("page2"),
			}, nil,
		).
		Once()

	mockClient.
		On(
			"ListServices",
			mock.Anything,
			mock.MatchedBy(func(input *vpclattice.ListServicesInput) bool {
				return input.NextToken != nil
			}),
		).
		Return(
			&vpclattice.ListServicesOutput{
				Items: []lattice.ServiceSummary{
					{Id: ptr.String("svc-100"), Name: ptr.String("service-100")},
				},
			}, nil,
		).
		Once()

	mockClient.
		On("ListListeners", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListListenersOutput{
				Items: []lattice.ListenerSummary{
					{Id: ptr.String("l-1"), Name: ptr.String("listener")},
				},
			}, nil,
		)

	mockClient.
		On("ListRules", mock.Anything, mock.Anything).
		Return(
			&vpclattice.ListRulesOutput{
				Items: []lattice.RuleSummary{
					{
						Id:        ptr.String("r-1"),
						Arn:       ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:rule/r-1"),
						Name:      ptr.String("rule"),
						Priority:  ptr.Int32(1),
						IsDefault: ptr.Bool(false),
					},
				},
			}, nil,
		)

	lister := &VPCLatticeListenerRuleLister{svc: mockClient}

	resources, err := lister.List(context.TODO(), testVPCLatticeListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 101)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListenerRule_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockVPCLatticeClient)

	rule := &VPCLatticeListenerRule{
		svc:        mockClient,
		ServiceID:  ptr.String("svc-1"),
		ListenerID: ptr.String("l-1"),
		ID:         ptr.String("r-1"),
	}

	mockClient.
		On(
			"DeleteRule",
			mock.Anything,
			&vpclattice.DeleteRuleInput{
				ServiceIdentifier:  rule.ServiceID,
				ListenerIdentifier: rule.ListenerID,
				RuleIdentifier:     rule.ID,
			},
		).
		Return(&vpclattice.DeleteRuleOutput{}, nil)

	err := rule.Remove(context.TODO())
	assertions.NoError(err)

	mockClient.AssertExpectations(t)
}

func Test_Mock_VPCLatticeListenerRule_Filter_Default(t *testing.T) {
	assertions := assert.New(t)

	rule := &VPCLatticeListenerRule{
		IsDefault: true,
		Name:      ptr.String("default"),
	}

	err := rule.Filter()
	assertions.NotNil(err)
	assertions.Contains(err.Error(), "cannot delete default listener rule")
}

func Test_Mock_VPCLatticeListenerRule_Filter_NonDefault(t *testing.T) {
	assertions := assert.New(t)

	rule := &VPCLatticeListenerRule{
		IsDefault: false,
		Name:      ptr.String("custom-rule"),
	}

	assertions.Nil(rule.Filter())
}

func Test_Mock_VPCLatticeListenerRule_Properties(t *testing.T) {
	assertions := assert.New(t)

	rule := VPCLatticeListenerRule{
		ID:           ptr.String("rule-12345"),
		ARN:          ptr.String("arn:aws:vpc-lattice:us-east-1:123456789012:rule/rule-12345"),
		Name:         ptr.String("my-rule"),
		ServiceID:    ptr.String("svc-111"),
		ServiceName:  ptr.String("my-service"),
		ListenerID:   ptr.String("listener-222"),
		ListenerName: ptr.String("my-listener"),
		Priority:     ptr.Int32(10),
		IsDefault:    false,
		Tags:         map[string]string{"Env": "prod"},
	}

	properties := rule.Properties()

	assertions.Equal("rule-12345", properties.Get("ID"))
	assertions.Equal("my-rule", properties.Get("Name"))
	assertions.Equal("my-service", properties.Get("ServiceName"))
	assertions.Equal("my-listener", properties.Get("ListenerName"))
	assertions.Equal("10", properties.Get("Priority"))
	assertions.Equal("prod", properties.Get("tag:Env"))
}

func Test_Mock_VPCLatticeListenerRule_String(t *testing.T) {
	assertions := assert.New(t)

	rule := VPCLatticeListenerRule{
		ServiceName:  ptr.String("my-service"),
		ListenerName: ptr.String("my-listener"),
		Name:         ptr.String("my-rule"),
	}

	assertions.Equal("my-service -> my-listener -> my-rule", rule.String())
}
