package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/networkmanager"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const NetworkManagerLinkResource = "NetworkManagerLink"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkManagerLinkResource,
		Scope:    nuke.Account,
		Resource: &NetworkManagerLink{},
		Lister:   &NetworkManagerLinkLister{},
		DependsOn: []string{
			NetworkManagerConnectionResource,
		},
	})
}

type NetworkManagerLinkLister struct {
	svc NetworkManagerClient
}

func (l *NetworkManagerLinkLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := l.svc
	if svc == nil {
		svc = networkmanager.NewFromConfig(*opts.Config)
	}

	var resources []resource.Resource

	// First, list all global networks
	globalNetworkIDs, err := listNetworkManagerGlobalNetworkIDs(ctx, svc)
	if err != nil {
		return nil, err
	}

	// For each global network, list all links
	for _, gnID := range globalNetworkIDs {
		linkParams := &networkmanager.GetLinksInput{
			GlobalNetworkId: &gnID,
		}
		for {
			resp, err := svc.GetLinks(ctx, linkParams)
			if err != nil {
				return nil, err
			}

			for _, link := range resp.Links {
				tags := make(map[string]string)
				for _, tag := range link.Tags {
					if tag.Key != nil && tag.Value != nil {
						tags[*tag.Key] = *tag.Value
					}
				}

				var state *string
				if link.State != "" {
					s := string(link.State)
					state = &s
				}

				resources = append(resources, &NetworkManagerLink{
					svc:             svc,
					ID:              link.LinkId,
					GlobalNetworkID: link.GlobalNetworkId,
					SiteID:          link.SiteId,
					Description:     link.Description,
					State:           state,
					Tags:            tags,
				})
			}

			if resp.NextToken == nil {
				break
			}
			linkParams.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

type NetworkManagerLink struct {
	svc             NetworkManagerClient
	ID              *string
	GlobalNetworkID *string
	SiteID          *string
	Description     *string
	State           *string
	Tags            map[string]string
}

func (r *NetworkManagerLink) Filter() error {
	if r.State != nil && strings.EqualFold(*r.State, "DELETING") {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *NetworkManagerLink) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteLink(ctx, &networkmanager.DeleteLinkInput{
		GlobalNetworkId: r.GlobalNetworkID,
		LinkId:          r.ID,
	})
	return err
}

func (r *NetworkManagerLink) Properties() types.Properties {
	properties := types.NewPropertiesFromStruct(r)
	for key, val := range r.Tags {
		properties.SetTag(&key, val)
	}
	return properties
}

func (r *NetworkManagerLink) String() string {
	return *r.ID
}
