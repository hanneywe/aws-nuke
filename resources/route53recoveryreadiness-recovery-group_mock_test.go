package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness"
	r53rrtypes "github.com/aws/aws-sdk-go-v2/service/route53recoveryreadiness/types"
)

func Test_Mock_Route53RecoveryReadinessRecoveryGroup_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryReadinessClient)
	mockClient.On("ListRecoveryGroups", mock.Anything, mock.Anything).
		Return(&route53recoveryreadiness.ListRecoveryGroupsOutput{
			RecoveryGroups: []r53rrtypes.RecoveryGroupOutput{
				{
					RecoveryGroupName: ptr.String("my-recovery-group"),
					RecoveryGroupArn:  ptr.String("arn:aws:route53-recovery-readiness::123456789012:recovery-group/my-recovery-group"),
				},
			},
		}, nil)
	lister := &Route53RecoveryReadinessRecoveryGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryReadinessListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	rg := resources[0].(*Route53RecoveryReadinessRecoveryGroup)
	a.Equal("my-recovery-group", *rg.RecoveryGroupName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryReadinessRecoveryGroup_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryReadinessClient)
	mockClient.On("ListRecoveryGroups", mock.Anything, mock.Anything).
		Return(&route53recoveryreadiness.ListRecoveryGroupsOutput{RecoveryGroups: []r53rrtypes.RecoveryGroupOutput{}}, nil)
	lister := &Route53RecoveryReadinessRecoveryGroupLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryReadinessListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryReadinessRecoveryGroup_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryReadinessClient)
	rg := &Route53RecoveryReadinessRecoveryGroup{svc: mockClient, RecoveryGroupName: ptr.String("my-recovery-group")}
	mockClient.On("DeleteRecoveryGroup", mock.Anything,
		&route53recoveryreadiness.DeleteRecoveryGroupInput{RecoveryGroupName: rg.RecoveryGroupName}).
		Return(&route53recoveryreadiness.DeleteRecoveryGroupOutput{}, nil)
	a.NoError(rg.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryReadinessRecoveryGroup_Properties(t *testing.T) {
	a := assert.New(t)
	rg := Route53RecoveryReadinessRecoveryGroup{
		RecoveryGroupName: ptr.String("my-recovery-group"),
		RecoveryGroupArn:  ptr.String("arn:aws:route53-recovery-readiness::123456789012:recovery-group/my-recovery-group"),
	}
	a.Equal("my-recovery-group", rg.Properties().Get("RecoveryGroupName"))
}

func Test_Mock_Route53RecoveryReadinessRecoveryGroup_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-recovery-group", (&Route53RecoveryReadinessRecoveryGroup{RecoveryGroupName: ptr.String("my-recovery-group")}).String())
}
