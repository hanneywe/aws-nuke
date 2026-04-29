package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/iot"
)

func Test_Mock_IoTDimension_List_One(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListDimensions", mock.Anything, mock.Anything).
		Return(&iot.ListDimensionsOutput{
			DimensionNames: []string{"my-dimension"},
		}, nil)

	lister := &IoTDimensionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 1)

	dimension := resources[0].(*IoTDimension)
	assertions.Equal("my-dimension", *dimension.Name)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTDimension_List_Empty(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTClient)

	mockClient.On("ListDimensions", mock.Anything, mock.Anything).
		Return(&iot.ListDimensionsOutput{
			DimensionNames: []string{},
		}, nil)

	lister := &IoTDimensionLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testIoTListerOpts)
	assertions.NoError(err)
	assertions.Len(resources, 0)

	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTDimension_Remove(t *testing.T) {
	assertions := assert.New(t)
	mockClient := new(mockIoTClient)

	dimension := &IoTDimension{
		svc:  mockClient,
		Name: ptr.String("my-dimension"),
	}

	mockClient.On("DeleteDimension", mock.Anything, &iot.DeleteDimensionInput{
		Name: dimension.Name,
	}).Return(&iot.DeleteDimensionOutput{}, nil)

	assertions.NoError(dimension.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_IoTDimension_Properties(t *testing.T) {
	assertions := assert.New(t)

	dimension := IoTDimension{
		Name: ptr.String("my-dimension"),
	}

	props := dimension.Properties()
	assertions.Equal("my-dimension", props.Get("Name"))
}

func Test_Mock_IoTDimension_String(t *testing.T) {
	assertions := assert.New(t)
	dimension := IoTDimension{Name: ptr.String("my-dimension")}
	assertions.Equal("my-dimension", dimension.String())
}
