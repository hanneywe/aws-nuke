package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/chimesdkvoice"
	chimesdkvoicetypes "github.com/aws/aws-sdk-go-v2/service/chimesdkvoice/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ChimeSDKVoicePhoneNumberResource = "ChimeSDKVoicePhoneNumber"

func init() {
	registry.Register(&registry.Registration{
		Name:     ChimeSDKVoicePhoneNumberResource,
		Scope:    nuke.Account,
		Resource: &ChimeSDKVoicePhoneNumber{},
		Lister:   &ChimeSDKVoicePhoneNumberLister{},
	})
}

type ChimeSDKVoicePhoneNumberLister struct {
	svc ChimeSDKVoiceClient
}

func (l *ChimeSDKVoicePhoneNumberLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = chimesdkvoice.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	paginator := chimesdkvoice.NewListPhoneNumbersPaginator(svc, &chimesdkvoice.ListPhoneNumbersInput{
		MaxResults: nil,
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for i := range page.PhoneNumbers {
			pn := &page.PhoneNumbers[i]
			resources = append(resources, &ChimeSDKVoicePhoneNumber{
				svc:              svc,
				PhoneNumberID:    pn.PhoneNumberId,
				E164PhoneNumber:  pn.E164PhoneNumber,
				Status:           pn.Status,
				ProductType:      pn.ProductType,
				CreatedTimestamp: pn.CreatedTimestamp,
			})
		}
	}

	return resources, nil
}

type ChimeSDKVoicePhoneNumber struct {
	svc              ChimeSDKVoiceClient
	PhoneNumberID    *string
	E164PhoneNumber  *string
	Status           chimesdkvoicetypes.PhoneNumberStatus
	ProductType      chimesdkvoicetypes.PhoneNumberProductType
	CreatedTimestamp *time.Time
}

func (r *ChimeSDKVoicePhoneNumber) Filter() error {
	if r.Status == chimesdkvoicetypes.PhoneNumberStatusDeleteInProgress {
		return fmt.Errorf("already being deleted")
	}
	if r.Status == chimesdkvoicetypes.PhoneNumberStatusCancelled {
		return fmt.Errorf("already canceled")
	}
	if r.Status == chimesdkvoicetypes.PhoneNumberStatusReleaseInProgress {
		return fmt.Errorf("release in progress")
	}
	return nil
}

func (r *ChimeSDKVoicePhoneNumber) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePhoneNumber(ctx, &chimesdkvoice.DeletePhoneNumberInput{
		PhoneNumberId: r.PhoneNumberID,
	})
	return err
}

func (r *ChimeSDKVoicePhoneNumber) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *ChimeSDKVoicePhoneNumber) String() string {
	if r.E164PhoneNumber != nil {
		return *r.E164PhoneNumber
	}
	return *r.PhoneNumberID
}
