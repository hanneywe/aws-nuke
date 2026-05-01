package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/cleanrooms"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const CleanRoomsPrivacyBudgetTemplateResource = "CleanRoomsPrivacyBudgetTemplate"

func init() {
	registry.Register(&registry.Registration{
		Name:     CleanRoomsPrivacyBudgetTemplateResource,
		Scope:    nuke.Account,
		Resource: &CleanRoomsPrivacyBudgetTemplate{},
		Lister:   &CleanRoomsPrivacyBudgetTemplateLister{},
	})
}

type CleanRoomsPrivacyBudgetTemplateLister struct {
	svc CleanRoomsClient
}

func (l *CleanRoomsPrivacyBudgetTemplateLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = cleanrooms.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	membershipPaginator := cleanrooms.NewListMembershipsPaginator(svc, &cleanrooms.ListMembershipsInput{})
	for membershipPaginator.HasMorePages() {
		membershipResp, err := membershipPaginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for i := range membershipResp.MembershipSummaries {
			membership := &membershipResp.MembershipSummaries[i]
			tmplPaginator := cleanrooms.NewListPrivacyBudgetTemplatesPaginator(svc, &cleanrooms.ListPrivacyBudgetTemplatesInput{
				MembershipIdentifier: membership.Id,
			})
			for tmplPaginator.HasMorePages() {
				tmplResp, err := tmplPaginator.NextPage(ctx)
				if err != nil {
					return nil, err
				}
				for _, tmpl := range tmplResp.PrivacyBudgetTemplateSummaries {
					resources = append(resources, &CleanRoomsPrivacyBudgetTemplate{
						svc:                             svc,
						MembershipIdentifier:            membership.Id,
						PrivacyBudgetTemplateIdentifier: tmpl.Id,
						CollaborationID:                 tmpl.CollaborationId,
					})
				}
			}
		}
	}

	return resources, nil
}

type CleanRoomsPrivacyBudgetTemplate struct {
	svc                             CleanRoomsClient
	MembershipIdentifier            *string
	PrivacyBudgetTemplateIdentifier *string
	CollaborationID                 *string
}

func (r *CleanRoomsPrivacyBudgetTemplate) Remove(ctx context.Context) error {
	_, err := r.svc.DeletePrivacyBudgetTemplate(ctx, &cleanrooms.DeletePrivacyBudgetTemplateInput{
		MembershipIdentifier:            r.MembershipIdentifier,
		PrivacyBudgetTemplateIdentifier: r.PrivacyBudgetTemplateIdentifier,
	})
	return err
}

func (r *CleanRoomsPrivacyBudgetTemplate) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *CleanRoomsPrivacyBudgetTemplate) String() string {
	return *r.PrivacyBudgetTemplateIdentifier
}
