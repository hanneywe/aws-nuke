package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/athena"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAthenaClient struct {
	mock.Mock
}

func (m *mockAthenaClient) ListCapacityReservations(
	ctx context.Context, params *athena.ListCapacityReservationsInput,
	_ ...func(*athena.Options),
) (*athena.ListCapacityReservationsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*athena.ListCapacityReservationsOutput), args.Error(1)
}

func (m *mockAthenaClient) CancelCapacityReservation(
	ctx context.Context, params *athena.CancelCapacityReservationInput,
	_ ...func(*athena.Options),
) (*athena.CancelCapacityReservationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*athena.CancelCapacityReservationOutput), args.Error(1)
}

func (m *mockAthenaClient) DeleteCapacityReservation(
	ctx context.Context, params *athena.DeleteCapacityReservationInput,
	_ ...func(*athena.Options),
) (*athena.DeleteCapacityReservationOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*athena.DeleteCapacityReservationOutput), args.Error(1)
}

var testAthenaListerOpts = &nuke.ListerOpts{}
