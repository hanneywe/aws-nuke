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

const NetworkManagerSiteResource = "NetworkManagerSite"

func init() {
	registry.Register(&registry.Registration{
		Name:     NetworkManagerSiteResource,
		Scope:    nuke.Account,
		Resource: &NetworkManagerSite{},
		Lister:   &NetworkManagerSiteLister{},
		DependsOn: []string{
			NetworkManagerLinkResource,
			NetworkManagerDeviceResource,
		},
	})
}

type NetworkManagerSiteLister struct {
	svc NetworkManagerClient
}

func (l *NetworkManagerSiteLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
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

	// For each global network, list all sites
	for _, gnID := range globalNetworkIDs {
		siteParams := &networkmanager.GetSitesInput{
			GlobalNetworkId: &gnID,
		}
		for {
			resp, err := svc.GetSites(ctx, siteParams)
			if err != nil {
				return nil, err
			}

			for _, site := range resp.Sites {
				tags := make(map[string]string)
				for _, tag := range site.Tags {
					if tag.Key != nil && tag.Value != nil {
						tags[*tag.Key] = *tag.Value
					}
				}

				var state *string
				if site.State != "" {
					s := string(site.State)
					state = &s
				}

				resources = append(resources, &NetworkManagerSite{
					svc:             svc,
					ID:              site.SiteId,
					GlobalNetworkID: site.GlobalNetworkId,
					Description:     site.Description,
					State:           state,
					Tags:            tags,
				})
			}

			if resp.NextToken == nil {
				break
			}
			siteParams.NextToken = resp.NextToken
		}
	}

	return resources, nil
}

type NetworkManagerSite struct {
	svc             NetworkManagerClient
	ID              *string
	GlobalNetworkID *string
	Description     *string
	State           *string
	Tags            map[string]string
}

func (r *NetworkManagerSite) Filter() error {
	if r.State != nil && strings.EqualFold(*r.State, "DELETING") {
		return fmt.Errorf("already deleting")
	}
	return nil
}

func (r *NetworkManagerSite) Remove(ctx context.Context) error {
	_, err := r.svc.DeleteSite(ctx, &networkmanager.DeleteSiteInput{
		GlobalNetworkId: r.GlobalNetworkID,
		SiteId:          r.ID,
	})
	return err
}

func (r *NetworkManagerSite) Properties() types.Properties {
	properties := types.NewPropertiesFromStruct(r)
	for key, val := range r.Tags {
		properties.SetTag(&key, val)
	}
	return properties
}

func (r *NetworkManagerSite) String() string {
	return *r.ID
}
