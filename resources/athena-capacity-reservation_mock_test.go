package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/athena"
	athenatypes "github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func Test_Mock_AthenaCapacityReservation_List_One(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAthenaClient)
	mockClient.On("ListCapacityReservations", mock.Anything, mock.Anything).
		Return(&athena.ListCapacityReservationsOutput{
			CapacityReservations: []athenatypes.CapacityReservation{
				{Name: ptr.String("my-reservation"), Status: athenatypes.CapacityReservationStatusActive},
			},
		}, nil)
	lister := &AthenaCapacityReservationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAthenaListerOpts)
	a.NoError(err)
	a.Len(resources, 1)
	a.Equal("my-reservation", *resources[0].(*AthenaCapacityReservation).Name)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AthenaCapacityReservation_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAthenaClient)
	mockClient.On("ListCapacityReservations", mock.Anything, mock.Anything).
		Return(&athena.ListCapacityReservationsOutput{CapacityReservations: []athenatypes.CapacityReservation{}}, nil)
	lister := &AthenaCapacityReservationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAthenaListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AthenaCapacityReservation_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAthenaClient)
	r := &AthenaCapacityReservation{
		svc:    mockClient,
		Name:   ptr.String("my-reservation"),
		Status: athenatypes.CapacityReservationStatusActive,
	}
	mockClient.On("CancelCapacityReservation", mock.Anything, &athena.CancelCapacityReservationInput{Name: r.Name}).
		Return(&athena.CancelCapacityReservationOutput{}, nil)
	mockClient.On("DeleteCapacityReservation", mock.Anything, &athena.DeleteCapacityReservationInput{Name: r.Name}).
		Return(&athena.DeleteCapacityReservationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	a.Equal(athenatypes.CapacityReservationStatusCancelled, r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AthenaCapacityReservation_Remove_SecondCall(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAthenaClient)
	r := &AthenaCapacityReservation{
		svc:    mockClient,
		Name:   ptr.String("my-reservation"),
		Status: athenatypes.CapacityReservationStatusActive,
	}
	mockClient.On("CancelCapacityReservation", mock.Anything, &athena.CancelCapacityReservationInput{Name: r.Name}).
		Return(&athena.CancelCapacityReservationOutput{}, nil)
	mockClient.On("DeleteCapacityReservation", mock.Anything, &athena.DeleteCapacityReservationInput{Name: r.Name}).
		Return(&athena.DeleteCapacityReservationOutput{}, nil)

	// First call: cancels then deletes
	a.NoError(r.Remove(context.TODO()))
	// Second call: should skip cancel, only delete
	a.NoError(r.Remove(context.TODO()))

	mockClient.AssertNumberOfCalls(t, "CancelCapacityReservation", 1)
	mockClient.AssertNumberOfCalls(t, "DeleteCapacityReservation", 2)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AthenaCapacityReservation_Remove_AlreadyCancelled(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAthenaClient)
	r := &AthenaCapacityReservation{
		svc:    mockClient,
		Name:   ptr.String("my-reservation"),
		Status: athenatypes.CapacityReservationStatusCancelled,
	}
	// Should NOT call CancelCapacityReservation since already canceled
	mockClient.On("DeleteCapacityReservation", mock.Anything, &athena.DeleteCapacityReservationInput{Name: r.Name}).
		Return(&athena.DeleteCapacityReservationOutput{}, nil)
	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AthenaCapacityReservation_Properties(t *testing.T) {
	a := assert.New(t)
	r := AthenaCapacityReservation{Name: ptr.String("my-reservation")}
	a.Equal("my-reservation", r.Properties().Get("Name"))
}

func Test_Mock_AthenaCapacityReservation_String(t *testing.T) {
	a := assert.New(t)
	a.Equal("my-reservation", (&AthenaCapacityReservation{Name: ptr.String("my-reservation")}).String())
}
