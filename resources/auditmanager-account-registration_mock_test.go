package resources

import (
	"context"
	"testing"

	"github.com/gotidy/ptr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/service/auditmanager"
	auditmanagertypes "github.com/aws/aws-sdk-go-v2/service/auditmanager/types"
)

func Test_Mock_AuditManagerAccountRegistration_List(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAuditmanagerClient)

	mockClient.On("GetAccountStatus", mock.Anything, mock.Anything).
		Return(&auditmanager.GetAccountStatusOutput{
			Status: auditmanagertypes.AccountStatusActive,
		}, nil)

	lister := &AuditManagerAccountRegistrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAuditmanagerListerOpts)
	a.NoError(err)
	a.Len(resources, 1)

	r := resources[0].(*AuditManagerAccountRegistration)
	a.Equal("ACTIVE", *r.Status)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AuditManagerAccountRegistration_List_Empty(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAuditmanagerClient)

	mockClient.On("GetAccountStatus", mock.Anything, mock.Anything).
		Return(&auditmanager.GetAccountStatusOutput{
			Status: auditmanagertypes.AccountStatusInactive,
		}, nil)

	lister := &AuditManagerAccountRegistrationLister{svc: mockClient}
	resources, err := lister.List(context.TODO(), testAuditmanagerListerOpts)
	a.NoError(err)
	a.Len(resources, 0)
	mockClient.AssertExpectations(t)
}

func Test_Mock_AuditManagerAccountRegistration_Remove(t *testing.T) {
	a := assert.New(t)
	mockClient := new(mockAuditmanagerClient)

	r := &AuditManagerAccountRegistration{
		svc:    mockClient,
		Status: ptr.String("ACTIVE"),
	}

	mockClient.On("DeregisterAccount", mock.Anything,
		&auditmanager.DeregisterAccountInput{}).
		Return(&auditmanager.DeregisterAccountOutput{}, nil)

	a.NoError(r.Remove(context.TODO()))
	mockClient.AssertExpectations(t)
}

func Test_Mock_AuditManagerAccountRegistration_Properties(t *testing.T) {
	a := assert.New(t)
	r := &AuditManagerAccountRegistration{
		Status: ptr.String("ACTIVE"),
	}
	props := r.Properties()
	a.Equal("ACTIVE", props.Get("Status"))
}

func Test_Mock_AuditManagerAccountRegistration_String(t *testing.T) {
	a := assert.New(t)
	r := &AuditManagerAccountRegistration{
		Status: ptr.String("ACTIVE"),
	}
	a.Equal("ACTIVE", r.String())
}
