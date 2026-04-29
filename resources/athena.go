package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/athena"
)

// AthenaClient is the interface for the Athena SDK v2 client methods used by sub-resources.
type AthenaClient interface {
	ListCapacityReservations(ctx context.Context, params *athena.ListCapacityReservationsInput,
		optFns ...func(*athena.Options)) (*athena.ListCapacityReservationsOutput, error)
	CancelCapacityReservation(ctx context.Context, params *athena.CancelCapacityReservationInput,
		optFns ...func(*athena.Options)) (*athena.CancelCapacityReservationOutput, error)
	DeleteCapacityReservation(ctx context.Context, params *athena.DeleteCapacityReservationInput,
		optFns ...func(*athena.Options)) (*athena.DeleteCapacityReservationOutput, error)
}
