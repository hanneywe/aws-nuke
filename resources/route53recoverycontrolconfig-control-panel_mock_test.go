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

func Test_Mock_Route53RecoveryControlConfigControlPanel_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	mockClient.On("ListControlPanels", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListControlPanelsOutput{
			ControlPanels: []r53rcctypes.ControlPanel{
				{
					ControlPanelArn:     ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel"),
					Name:                ptr.String("my-panel"),
					DefaultControlPanel: ptr.Bool(false),
				},
			},
		}, nil)
	lister := &Route53RecoveryControlConfigControlPanelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	cp := resources[0].(*Route53RecoveryControlConfigControlPanel)
	a.Equal("my-panel", *cp.Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigControlPanel_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	mockClient.On("ListControlPanels", mock.Anything, mock.Anything).
		Return(&route53recoverycontrolconfig.ListControlPanelsOutput{ControlPanels: []r53rcctypes.ControlPanel{}}, nil)
	lister := &Route53RecoveryControlConfigControlPanelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryControlConfigListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigControlPanel_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryControlConfigClient)
	panelArn := "arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel"
	cp := &Route53RecoveryControlConfigControlPanel{
		svc:             mockClient,
		ControlPanelArn: ptr.String(panelArn),
	}
	mockClient.On("DeleteControlPanel", mock.Anything,
		&route53recoverycontrolconfig.DeleteControlPanelInput{ControlPanelArn: cp.ControlPanelArn}).
		Return(&route53recoverycontrolconfig.DeleteControlPanelOutput{}, nil)
	a.NoError(cp.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryControlConfigControlPanel_Filter_Default(t *testing.T) {
	a := assert.New(t)
	cp := Route53RecoveryControlConfigControlPanel{DefaultControlPanel: ptr.Bool(true), Name: ptr.String("DefaultControlPanel")}
	a.Error(cp.Filter())
}

func Test_Mock_Route53RecoveryControlConfigControlPanel_Filter_NonDefault(t *testing.T) {
	a := assert.New(t)
	cp := Route53RecoveryControlConfigControlPanel{DefaultControlPanel: ptr.Bool(false), Name: ptr.String("my-panel")}
	a.NoError(cp.Filter())
}

func Test_Mock_Route53RecoveryControlConfigControlPanel_Properties(t *testing.T) {
	a := assert.New(t)
	cp := Route53RecoveryControlConfigControlPanel{
		ControlPanelArn: ptr.String("arn:aws:route53-recovery-control::123456789012:controlpanel/my-panel"),
		Name:            ptr.String("my-panel"),
	}
	a.Equal("my-panel", cp.Properties().Get("Name"))
}

func Test_Mock_Route53RecoveryControlConfigControlPanel_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-panel", (&Route53RecoveryControlConfigControlPanel{Name: ptr.String("my-panel")}).String())
}
