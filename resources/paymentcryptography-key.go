package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/paymentcryptography"
	paymentcryptographytypes "github.com/aws/aws-sdk-go-v2/service/paymentcryptography/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PaymentCryptographyKeyResource = "PaymentCryptographyKey"

func init() {
	registry.Register(&registry.Registration{
		Name:     PaymentCryptographyKeyResource,
		Scope:    nuke.Account,
		Resource: &PaymentCryptographyKey{},
		Lister:   &PaymentCryptographyKeyLister{},
	})
}

type PaymentCryptographyKeyLister struct {
	svc PaymentCryptographyClient
}

func (l *PaymentCryptographyKeyLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = paymentcryptography.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := paymentcryptography.NewListKeysPaginator(svc, &paymentcryptography.ListKeysInput{
		MaxResults: aws.Int32(100),
	})
	for paginator.HasMorePages() {
		listKeysOutput, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, keySummary := range listKeysOutput.Keys {
			resources = append(resources, &PaymentCryptographyKey{
				svc:      svc,
				KeyArn:   keySummary.KeyArn,
				KeyState: keySummary.KeyState,
			})
		}
	}

	return resources, nil
}

type PaymentCryptographyKey struct {
	svc      PaymentCryptographyClient
	KeyArn   *string
	KeyState paymentcryptographytypes.KeyState `property:"-"`
}

func (r *PaymentCryptographyKey) Filter() error {
	if r.KeyState == paymentcryptographytypes.KeyStateDeletePending ||
		r.KeyState == paymentcryptographytypes.KeyStateDeleteComplete {
		return fmt.Errorf("already %s", r.KeyState)
	}
	return nil
}

func (r *PaymentCryptographyKey) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteKey(ctx, &paymentcryptography.DeleteKeyInput{
		KeyIdentifier:   r.KeyArn,
		DeleteKeyInDays: aws.Int32(3),
	})
	return err
}

func (r *PaymentCryptographyKey) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PaymentCryptographyKey) String() string {
	return *r.KeyArn
}
