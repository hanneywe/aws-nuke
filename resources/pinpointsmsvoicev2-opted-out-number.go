package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const PinpointSMSVoiceV2OptedOutNumberResource = "PinpointSMSVoiceV2OptedOutNumber"

func init() {
	registry.Register(&registry.Registration{
		Name:     PinpointSMSVoiceV2OptedOutNumberResource,
		Scope:    nuke.Account,
		Resource: &PinpointSMSVoiceV2OptedOutNumber{},
		Lister:   &PinpointSMSVoiceV2OptedOutNumberLister{},
	})
}

type PinpointSMSVoiceV2OptedOutNumberLister struct {
	svc PinpointSMSVoiceV2Client
}

func (l *PinpointSMSVoiceV2OptedOutNumberLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = pinpointsmsvoicev2.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	listPaginator := pinpointsmsvoicev2.NewDescribeOptOutListsPaginator(svc, &pinpointsmsvoicev2.DescribeOptOutListsInput{})

	for listPaginator.HasMorePages() {
		listResp, err := listPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, optOutList := range listResp.OptOutLists {
			numberPaginator := pinpointsmsvoicev2.NewDescribeOptedOutNumbersPaginator(svc, &pinpointsmsvoicev2.DescribeOptedOutNumbersInput{
				OptOutListName: optOutList.OptOutListName,
			})

			for numberPaginator.HasMorePages() {
				numberResp, err := numberPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}

				for _, number := range numberResp.OptedOutNumbers {
					resources = append(resources, &PinpointSMSVoiceV2OptedOutNumber{
						svc:               svc,
						OptedOutNumber:    number.OptedOutNumber,
						OptOutListName:    optOutList.OptOutListName,
						OptedOutTimestamp: number.OptedOutTimestamp,
						EndUserOptedOut:   &number.EndUserOptedOut,
					})
				}
			}
		}
	}

	return resources, nil
}

type PinpointSMSVoiceV2OptedOutNumber struct {
	svc               PinpointSMSVoiceV2Client
	OptedOutNumber    *string
	OptOutListName    *string
	OptedOutTimestamp *time.Time
	EndUserOptedOut   *bool
}

func (r *PinpointSMSVoiceV2OptedOutNumber) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteOptedOutNumber(ctx, &pinpointsmsvoicev2.DeleteOptedOutNumberInput{
		OptOutListName: r.OptOutListName,
		OptedOutNumber: r.OptedOutNumber,
	})
	return err
}

func (r *PinpointSMSVoiceV2OptedOutNumber) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *PinpointSMSVoiceV2OptedOutNumber) String() string {
	return fmt.Sprintf("%s -> %s", *r.OptOutListName, *r.OptedOutNumber)
}
