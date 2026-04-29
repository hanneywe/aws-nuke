package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/location"
	locationtypes "github.com/aws/aws-sdk-go-v2/service/location/types"
)

func Test_Mock_LocationServiceRouteCalculator_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListRouteCalculators", mock.Anything, mock.Anything).
		Return(&location.ListRouteCalculatorsOutput{
			Entries: []locationtypes.ListRouteCalculatorsResponseEntry{
				{CalculatorName: ptr.String("my-calculator")},
			},
		}, nil)

	lister := &LocationServiceRouteCalculatorLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	rc := resources[0].(*LocationServiceRouteCalculator)
	a.Equal("my-calculator", *rc.CalculatorName)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceRouteCalculator_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	mockClient.On("ListRouteCalculators", mock.Anything, mock.Anything).
		Return(&location.ListRouteCalculatorsOutput{
			Entries: []locationtypes.ListRouteCalculatorsResponseEntry{},
		}, nil)

	lister := &LocationServiceRouteCalculatorLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testLocationServiceListerOpts)
	a.NoError(err)
	a.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceRouteCalculator_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockLocationServiceClient)

	rc := &LocationServiceRouteCalculator{
		svc:            mockClient,
		CalculatorName: ptr.String("my-calculator"),
	}

	mockClient.On("DeleteRouteCalculator", mock.Anything, &location.DeleteRouteCalculatorInput{
		CalculatorName: rc.CalculatorName,
	}).Return(&location.DeleteRouteCalculatorOutput{}, nil)

	a.NoError(rc.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_LocationServiceRouteCalculator_Properties(t *testing.T) {
	a := assert.New(t)
	rc := LocationServiceRouteCalculator{CalculatorName: ptr.String("my-calculator")}
	props := rc.Properties()
	a.Equal("my-calculator", props.Get("CalculatorName"))
}

func Test_Mock_LocationServiceRouteCalculator_String(t *testing.T) {
	a := assert.New(t)
	rc := LocationServiceRouteCalculator{CalculatorName: ptr.String("my-calculator")}
	a.Equal("my-calculator", rc.String())
}
