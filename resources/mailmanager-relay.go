package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/mailmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const MailManagerRelayResource = "MailManagerRelay"

func init() {
	registry.Register(&registry.Registration{
		Name:     MailManagerRelayResource,
		Scope:    nuke.Account,
		Resource: &MailManagerRelay{},
		Lister:   &MailManagerRelayLister{},
		DependsOn: []string{
			MailManagerRuleSetResource,
		},
	})
}

type MailManagerRelayLister struct {
	svc MailManagerClient
}

func (l *MailManagerRelayLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = mailmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &mailmanager.ListRelaysInput{}

	for {
		output, err := svc.ListRelays(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, relay := range output.Relays {
			resources = append(resources, &MailManagerRelay{
				svc:       svc,
				RelayID:   relay.RelayId,
				RelayName: relay.RelayName,
			})
		}

		if output.NextToken == nil {
			break
		}

		params.NextToken = output.NextToken
	}

	return resources, nil
}

type MailManagerRelay struct {
	svc       MailManagerClient
	RelayID   *string `property:"name=RelayId"`
	RelayName *string
}

func (r *MailManagerRelay) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteRelay(ctx, &mailmanager.DeleteRelayInput{
		RelayId: r.RelayID,
	})
	return err
}

func (r *MailManagerRelay) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *MailManagerRelay) String() string {
	return *r.RelayName
}
