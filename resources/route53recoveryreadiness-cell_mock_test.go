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

func Test_Mock_Route53RecoveryReadinessCell_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryReadinessClient)
	mockClient.On("ListCells", mock.Anything, mock.Anything).
		Return(&route53recoveryreadiness.ListCellsOutput{
			Cells: []r53rrtypes.CellOutput{
				{CellName: ptr.String("my-cell"), CellArn: ptr.String("arn:aws:route53-recovery-readiness::123456789012:cell/my-cell")},
			},
		}, nil)
	lister := &Route53RecoveryReadinessCellLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryReadinessListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	cell := resources[0].(*Route53RecoveryReadinessCell)
	a.Equal("my-cell", *cell.CellName)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryReadinessCell_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryReadinessClient)
	mockClient.On("ListCells", mock.Anything, mock.Anything).
		Return(&route53recoveryreadiness.ListCellsOutput{Cells: []r53rrtypes.CellOutput{}}, nil)
	lister := &Route53RecoveryReadinessCellLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testRoute53RecoveryReadinessListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryReadinessCell_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockRoute53RecoveryReadinessClient)
	cell := &Route53RecoveryReadinessCell{svc: mockClient, CellName: ptr.String("my-cell")}
	mockClient.On("DeleteCell", mock.Anything, &route53recoveryreadiness.DeleteCellInput{CellName: cell.CellName}).
		Return(&route53recoveryreadiness.DeleteCellOutput{}, nil)
	a.NoError(cell.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_Route53RecoveryReadinessCell_Properties(t *testing.T) {
	a := assert.New(t)
	cell := Route53RecoveryReadinessCell{
		CellName: ptr.String("my-cell"),
		CellArn:  ptr.String("arn:aws:route53-recovery-readiness::123456789012:cell/my-cell"),
	}
	a.Equal("my-cell", cell.Properties().Get("CellName"))
}

func Test_Mock_Route53RecoveryReadinessCell_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-cell", (&Route53RecoveryReadinessCell{CellName: ptr.String("my-cell")}).String())
}
