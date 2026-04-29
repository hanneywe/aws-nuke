package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/comprehend"
	comprehendtypes "github.com/aws/aws-sdk-go-v2/service/comprehend/types"
)

func Test_Mock_ComprehendFlywheel_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockComprehendClient)
	mockClient.On("ListFlywheels", mock.Anything, mock.Anything).
		Return(&comprehend.ListFlywheelsOutput{
			FlywheelSummaryList: []comprehendtypes.FlywheelSummary{
				{
					FlywheelArn: ptr.String("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel"),
					ModelType:   comprehendtypes.ModelTypeDocumentClassifier,
					Status:      comprehendtypes.FlywheelStatusActive,
				},
			},
		}, nil)
	lister := &ComprehendFlywheelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testComprehendListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	fw := resources[0].(*ComprehendFlywheel)
	a.Equal("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel", *fw.FlywheelArn)
	a.Equal(comprehendtypes.ModelTypeDocumentClassifier, fw.ModelType)
	a.Equal(comprehendtypes.FlywheelStatusActive, fw.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ComprehendFlywheel_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockComprehendClient)
	mockClient.On("ListFlywheels", mock.Anything, mock.Anything).
		Return(&comprehend.ListFlywheelsOutput{
			FlywheelSummaryList: []comprehendtypes.FlywheelSummary{},
		}, nil)
	lister := &ComprehendFlywheelLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testComprehendListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_ComprehendFlywheel_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockComprehendClient)
	fw := &ComprehendFlywheel{
		svc:         mockClient,
		FlywheelArn: ptr.String("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel"),
	}
	mockClient.On("DeleteFlywheel", mock.Anything, &comprehend.DeleteFlywheelInput{
		FlywheelArn: fw.FlywheelArn,
	}).Return(&comprehend.DeleteFlywheelOutput{}, nil)
	a.NoError(fw.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_ComprehendFlywheel_Properties(t *testing.T) {
	a := assert.New(t)
	fw := ComprehendFlywheel{
		FlywheelArn: ptr.String("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel"),
		ModelType:   comprehendtypes.ModelTypeDocumentClassifier,
		Status:      comprehendtypes.FlywheelStatusActive,
	}
	props := fw.Properties()
	a.Equal("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel", props.Get("FlywheelArn"))
	a.Equal(string(comprehendtypes.ModelTypeDocumentClassifier), props.Get("ModelType"))
	a.Equal(string(comprehendtypes.FlywheelStatusActive), props.Get("Status"))
}

func Test_Mock_ComprehendFlywheel_String(t *testing.T) {
	a := assert.New(t)
	fw := ComprehendFlywheel{
		FlywheelArn: ptr.String("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel"),
	}
	a.Equal("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel", fw.String())
}

func Test_Mock_ComprehendFlywheel_Filter_Deleting(t *testing.T) {
	a := assert.New(t)
	fw := ComprehendFlywheel{
		FlywheelArn: ptr.String("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel"),
		Status:      comprehendtypes.FlywheelStatusDeleting,
	}
	a.Error(fw.Filter())
}

func Test_Mock_ComprehendFlywheel_Filter_Active(t *testing.T) {
	a := assert.New(t)
	fw := ComprehendFlywheel{
		FlywheelArn: ptr.String("arn:aws:comprehend:us-east-1:123456789012:flywheel/my-flywheel"),
		Status:      comprehendtypes.FlywheelStatusActive,
	}
	a.NoError(fw.Filter())
}
