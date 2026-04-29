package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	auditmanagertypes "github.com/aws/aws-sdk-go-v2/service/auditmanager/types"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const AuditManagerAccountRegistrationResource = "AuditManagerAccountRegistration"

func init() {
	registry.Register(&registry.Registration{
		Name:     AuditManagerAccountRegistrationResource,
		Scope:    nuke.Account,
		Resource: &AuditManagerAccountRegistration{},
		Lister:   &AuditManagerAccountRegistrationLister{},
	})
}

type AuditManagerAccountRegistrationLister struct {
	svc AuditmanagerClient
}

func (l *AuditManagerAccountRegistrationLister) List(ctx context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)
	svc := l.svc
	if svc == nil {
		svc = auditmanager.NewFromConfig(*opts.Config)
	}
	var resources []resource.Resource

	resp, err := svc.GetAccountStatus(ctx, &auditmanager.GetAccountStatusInput{})
	if err != nil {
		return nil, err
	}

	if resp.Status == auditmanagertypes.AccountStatusActive {
		status := string(resp.Status)
		resources = append(resources, &AuditManagerAccountRegistration{
			svc:    svc,
			Status: &status,
		})
	}

	return resources, nil
}

type AuditManagerAccountRegistration struct {
	svc    AuditmanagerClient
	Status *string
}

func (r *AuditManagerAccountRegistration) Remove(ctx context.Context) error {
	_, err := r.svc.DeregisterAccount(ctx, &auditmanager.DeregisterAccountInput{})
	return err
}

func (r *AuditManagerAccountRegistration) Properties() types.Properties {
	return types.NewPropertiesFromStruct(r)
}

func (r *AuditManagerAccountRegistration) String() string {
	return *r.Status
}
