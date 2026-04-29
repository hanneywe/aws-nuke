package resources

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const IoTAuditSuppressionResource = "IoTAuditSuppression"

func init() {
	registry.Register(&registry.Registration{
		Name:     IoTAuditSuppressionResource,
		Scope:    nuke.Account,
		Resource: &IoTAuditSuppression{},
		Lister:   &IoTAuditSuppressionLister{},
	})
}

type IoTAuditSuppressionLister struct {
	svc IoTClient
}

func (l *IoTAuditSuppressionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = iot.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource
	params := &iot.ListAuditSuppressionsInput{
		MaxResults: aws.Int32(100),
	}

	for {
		resp, err := svc.ListAuditSuppressions(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, s := range resp.Suppressions {
			var resourceIdentifierAccount *string
			if s.ResourceIdentifier != nil && s.ResourceIdentifier.Account != nil {
				resourceIdentifierAccount = s.ResourceIdentifier.Account
			}

			resources = append(resources, &IoTAuditSuppression{
				svc:                       svc,
				CheckName:                 s.CheckName,
				ResourceIdentifier:        s.ResourceIdentifier,
				Description:               s.Description,
				ExpirationDate:            s.ExpirationDate,
				SuppressIndefinitely:      s.SuppressIndefinitely,
				ResourceIdentifierAccount: resourceIdentifierAccount,
			})
		}

		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type IoTAuditSuppression struct {
	svc                       IoTClient
	CheckName                 *string
	ResourceIdentifier        *iottypes.ResourceIdentifier `property:"-"`
	Description               *string
	ExpirationDate            *time.Time
	SuppressIndefinitely      *bool
	ResourceIdentifierAccount *string
}

func (r *IoTAuditSuppression) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteAuditSuppression(ctx, &iot.DeleteAuditSuppressionInput{
		CheckName:          r.CheckName,
		ResourceIdentifier: r.ResourceIdentifier,
	})
	return err
}

func (r *IoTAuditSuppression) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *IoTAuditSuppression) String() string {
	return *r.CheckName
}
