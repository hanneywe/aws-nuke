package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EventBridgeEventBusResource = "EventBridgeEventBus"

func init() {
	registry.Register(&registry.Registration{
		Name:     EventBridgeEventBusResource,
		Scope:    nuke.Account,
		Resource: &EventBridgeEventBus{},
		Lister:   &EventBridgeEventBusLister{},
	})
}

type EventBridgeEventBusLister struct {
	svc EventBridgeClient
}

func (l *EventBridgeEventBusLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = eventbridge.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	params := &eventbridge.ListEventBusesInput{}

	for {
		resp, err := svc.ListEventBuses(ctx, params)
		if err != nil {
			return nil, err
		}

		for i := range resp.EventBuses {
			item := &resp.EventBuses[i]
			resources = append(resources, &EventBridgeEventBus{
				svc:  svc,
				Name: item.Name,
				Arn:  item.Arn,
			})
		}
		if resp.NextToken == nil {
			break
		}
		params.NextToken = resp.NextToken
	}

	return resources, nil
}

type EventBridgeEventBus struct {
	svc  EventBridgeClient
	Name *string
	Arn  *string
}

func (r *EventBridgeEventBus) Filter() error {
	if *r.Name == "default" {
		return fmt.Errorf("cannot delete default event bus")
	}
	return nil
}

func (r *EventBridgeEventBus) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteEventBus(ctx, &eventbridge.DeleteEventBusInput{
		Name: r.Name,
	})
	return err
}

func (r *EventBridgeEventBus) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EventBridgeEventBus) String() string {
	return *r.Name
}
