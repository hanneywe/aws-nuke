package resources

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/auditmanager"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

type mockAuditmanagerClient struct {
	mock.Mock
}

func (m *mockAuditmanagerClient) GetAccountStatus(
	ctx context.Context, params *auditmanager.GetAccountStatusInput,
	_ ...func(*auditmanager.Options),
) (*auditmanager.GetAccountStatusOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*auditmanager.GetAccountStatusOutput), args.Error(1)
}

func (m *mockAuditmanagerClient) DeregisterAccount(
	ctx context.Context, params *auditmanager.DeregisterAccountInput,
	_ ...func(*auditmanager.Options),
) (*auditmanager.DeregisterAccountOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*auditmanager.DeregisterAccountOutput), args.Error(1)
}

var testAuditmanagerListerOpts = &nuke.ListerOpts{}
