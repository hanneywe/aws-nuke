package resources

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/organizations"
	orgtypes "github.com/aws/aws-sdk-go-v2/service/organizations/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	libsettings "github.com/ekristen/libnuke/pkg/settings"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const OrganizationsAccountResource = "OrganizationsAccount"

func init() {
	registry.Register(&registry.Registration{
		Name:     OrganizationsAccountResource,
		Scope:    nuke.Account,
		Resource: &OrganizationsAccount{},
		Lister:   &OrganizationsAccountLister{},
		Settings: []string{"CloseAccount"},
	})
}

type OrganizationsAccountLister struct {
	svc OrganizationsClient
}

func (l *OrganizationsAccountLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = organizations.NewFromConfig(*opts.Config)
	}

	descResp, err := svc.DescribeOrganization(ctx, &organizations.DescribeOrganizationInput{})
	if err != nil {
		return nil, err
	}

	var managementAccountID string
	if descResp.Organization != nil && descResp.Organization.MasterAccountId != nil {
		managementAccountID = *descResp.Organization.MasterAccountId
	}

	var resources []resource.Resource

	paginator := organizations.NewListAccountsPaginator(svc, &organizations.ListAccountsInput{})
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, account := range resp.Accounts {
			isManagement := account.Id != nil && *account.Id == managementAccountID
			resources = append(resources, &OrganizationsAccount{
				svc:               svc,
				ID:                account.Id,
				Name:              account.Name,
				Email:             account.Email,
				Status:            account.Status,
				ManagementAccount: &isManagement,
			})
		}
	}

	return resources, nil
}

type OrganizationsAccount struct {
	svc               OrganizationsClient
	ID                *string
	Name              *string
	Email             *string
	Status            orgtypes.AccountStatus
	ManagementAccount *bool
	settings          *libsettings.Setting
}

func (r *OrganizationsAccount) Settings(setting *libsettings.Setting) {
	r.settings = setting
}

func (r *OrganizationsAccount) Filter() error {
	if r.Status == orgtypes.AccountStatusSuspended {
		return fmt.Errorf("account is SUSPENDED")
	}
	if r.Status == orgtypes.AccountStatusPendingClosure {
		return fmt.Errorf("account is PENDING_CLOSURE")
	}
	if r.ManagementAccount != nil && *r.ManagementAccount {
		return fmt.Errorf("cannot close management account")
	}
	return nil
}

func (r *OrganizationsAccount) Remove(ctx context.Context) error {
	if !r.settings.GetBool("CloseAccount") {
		return fmt.Errorf("CloseAccount setting must be enabled to close account %s", *r.ID)
	}
	_, err := r.svc.CloseAccount(ctx, &organizations.CloseAccountInput{
		AccountId: r.ID,
	})
	return err
}

func (r *OrganizationsAccount) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *OrganizationsAccount) String() string {
	return *r.ID
}
