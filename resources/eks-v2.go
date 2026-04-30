package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

// EKSv2Client is the interface for the EKS SDK v2 client methods used by new resources.
type EKSv2Client interface {
	ListEksAnywhereSubscriptions(ctx context.Context, params *eks.ListEksAnywhereSubscriptionsInput,
		optFns ...func(*eks.Options)) (*eks.ListEksAnywhereSubscriptionsOutput, error)
	DeleteEksAnywhereSubscription(ctx context.Context, params *eks.DeleteEksAnywhereSubscriptionInput,
		optFns ...func(*eks.Options)) (*eks.DeleteEksAnywhereSubscriptionOutput, error)
}
