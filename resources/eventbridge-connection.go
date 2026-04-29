package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventbridgetypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const EventBridgeConnectionResource = "EventBridgeConnection"

func init() {
	registry.Register(&registry.Registration{
		Name:     EventBridgeConnectionResource,
		Scope:    nuke.Account,
		Resource: &EventBridgeConnection{},
		Lister:   &EventBridgeConnectionLister{},
	})
}

type EventBridgeConnectionLister struct {
	svc EventBridgeClient
}

func (l *EventBridgeConnectionLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = eventbridge.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	params := &eventbridge.ListConnectionsInput{}

	for {
		output, err := svc.ListConnections(ctx, params)
		if err != nil {
			return nil, err
		}

		for _, connection := range output.Connections {
			connectionState := string(connection.ConnectionState)
			resources = append(resources, &EventBridgeConnection{
				svc:             svc,
				Name:            connection.Name,
				ConnectionArn:   connection.ConnectionArn,
				ConnectionState: &connectionState,
			})
		}

		if output.NextToken == nil {
			break
		}
		params.NextToken = output.NextToken
	}

	return resources, nil
}

type EventBridgeConnection struct {
	svc             EventBridgeClient
	Name            *string
	ConnectionArn   *string
	ConnectionState *string
}

func (r *EventBridgeConnection) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteConnection(ctx, &eventbridge.DeleteConnectionInput{
		Name: r.Name,
	})
	return err
}

func (r *EventBridgeConnection) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *EventBridgeConnection) String() string {
	return *r.Name
}

func (r *EventBridgeConnection) Filter() error {
	if r.ConnectionState != nil {
		if eventbridgetypes.ConnectionState(*r.ConnectionState) == eventbridgetypes.ConnectionStateDeleting {
			return fmt.Errorf("already %s", *r.ConnectionState)
		}
	}
	return nil
}
