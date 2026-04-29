package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig"
	r53rcctypes "github.com/aws/aws-sdk-go-v2/service/route53recoverycontrolconfig/types"
)

func Test_Mock_Route53RecoveryControlConfigRoutingControl_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)

	// First call: list control panels
	mockClient.On("ListControlPanels", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListControlPanelsOutput{
			ControlPanels: []r53rcctypes.ControlPanel{
				{ControlPanelArn: ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel")},
			},
		}, nil)

	// Second call: list routing controls for the panel
	mockClient.On("ListRoutingControls", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListRoutingControlsOutput{
			RoutingControls: []r53rcctypes.RoutingControl{
				{
					RoutingControlArn: ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel/routingcontrol/my-rc"),
					Name:              ptr.String("my-rc"),
					ControlPanelArn:   ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel"),
				},
			},
		}, nil)

	lister := &Route53RecoveryControlConfigRoutingControlLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	rc := resources[0].(*Route53RecoveryControlConfigRoutingControl)
	a.Equal("my-rc", *rc.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigRoutingControl_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	mockClient.On("ListControlPanels", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListControlPanelsOutput{ControlPanels: []r53rcctypes.ControlPanel{}}, nil)
	lister := &Route53RecoveryControlConfigRoutingControlLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigRoutingControl_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	rc := &Route53RecoveryControlConfigRoutingControl{
		svc:               mockClient,
		RoutingControlArn: ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel/routingcontrol/my-rc"),
	}
	mockClient.On("DeleteRoutingControl", mock.Anything, &route53recoverycontrolconfig.DeleteRoutingControlInput{
		RoutingControlArn: rc.RoutingControlArn,
	}).Return(&route53recoverycontrolconfig.DeleteRoutingControlOutput{}, nil)
	a.NoError(rc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigRoutingControl_Properties(t *testing.T) {
	a := assert.New(t)
	rc := Route53RecoveryControlConfigRoutingControl{
		RoutingControlArn: ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel/routingcontrol/my-rc"),
		Name:              ptr.String("my-rc"),
		ControlPanelArn:   ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel"),
	}
	a.Equal("my-rc", rc.Properties().Get("Name"))
}

func Test_Mock_Route53RecoveryControlConfigRoutingControl_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-rc", (&Route53RecoveryControlConfigRoutingControl{Name: ptr.String("my-rc")}).String())
}
