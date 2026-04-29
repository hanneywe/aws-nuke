package resources

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
)

// AuditmanagerClient is the interface for the auditmanager SDK client methods.
type AuditmanagerClient interface {
	GetAccountStatus(ctx context.Context, params *auditmanager.GetAccountStatusInput,
		optFns ...func(*auditmanager.Options)) (*auditmanager.GetAccountStatusOutput, error)
	DeregisterAccount(ctx context.Context, params *auditmanager.DeregisterAccountInput,
		optFns ...func(*auditmanager.Options)) (*auditmanager.DeregisterAccountOutput, error)
}
